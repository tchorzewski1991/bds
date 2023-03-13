package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/tchorzewski1991/bds/app/services/books-api/handlers"
	"github.com/tchorzewski1991/bds/business/sys/database"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/ardanlabs/conf/v3"
	"github.com/emadolsky/automaxprocs/maxprocs"
	"github.com/tchorzewski1991/bds/base/logger"
	_ "go.uber.org/automaxprocs"
	"go.uber.org/zap"
)

var build = "develop"
var service = "BOOKS-API"

func main() {
	// ================================================================================================================
	// Set GOMAXPROCS

	// Set the correct number of threads for the service based on what is
	// available either by the machine or quotas.
	if _, err := maxprocs.Set(); err != nil {
		fmt.Println("Setting maxprocs error: %w", err)
	}
	// TODO: Figure out how and where to put info about cpu.
	_ = runtime.GOMAXPROCS(0)

	// ================================================================================================================
	// Configure app logger

	fields := []logger.Field{
		{
			Name: "service", Value: service,
		},
	}
	l, err := logger.New(fields...)
	if err != nil {
		fmt.Println("Constructing logger error:", err)
		os.Exit(1)
	}
	defer func() { _ = l.Sync() }()

	// ================================================================================================================
	// Run app

	if err = run(l); err != nil {
		l.Errorf("Running app error: %s", err)
		os.Exit(1)
	}
}

func run(logger *zap.SugaredLogger) error {
	// ================================================================================================================
	// Configuration

	cfg := struct {
		conf.Version
		Api struct {
			Host            string        `conf:"default:0.0.0.0:3000"`
			DebugHost       string        `conf:"default:0.0.0.0:4000"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:10s"`
			IdleTimeout     time.Duration `conf:"default:30s"`
			ShutdownTimeout time.Duration `conf:"default:30s"`
		}
		DB struct {
			User string `conf:"default:postgres"`
			Pass string `conf:"default:password,mask"`
			Host string `conf:"default:localhost"`
			Name string `conf:"default:bds"`
		}
	}{
		Version: conf.Version{
			Build: build,
			Desc:  "Current build version",
		},
	}

	const prefix = "BOOKS"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config output: %w", err)
	}
	logger.Infow("Config parsed", "config", out)

	// ================================================================================================================
	// Database support

	logger.Info("Starting database")

	db, err := database.Open(database.Config{
		User: cfg.DB.User,
		Pass: cfg.DB.Pass,
		Host: cfg.DB.Host,
		Name: cfg.DB.Name,
	})
	if err != nil {
		return fmt.Errorf("connecting to db: %w", err)
	}
	defer func() {
		logger.Infow("Database shutdown")
		_ = db.Close()
	}()

	// ================================================================================================================
	// Start Debug service

	logger.Infow("Starting debug service", "host", cfg.Api.DebugHost)

	debugMux := handlers.DebugMux(handlers.DebugMuxConfig{
		Build:  build,
		Logger: logger,
		DB:     db,
	})

	go func() {
		if err = http.ListenAndServe(cfg.Api.DebugHost, debugMux); err != nil {
			logger.Errorw("Debug service shutdown", "host", cfg.Api.DebugHost, "error", err)
		}
	}()

	// ================================================================================================================
	// Starting App

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)

	apiMux := handlers.ApiMux(handlers.ApiMuxConfig{
		Shutdown: shutdown,
		Logger:   logger,
		DB:       db,
	})

	apiSrv := http.Server{
		Addr:         cfg.Api.Host,
		Handler:      apiMux,
		ReadTimeout:  cfg.Api.ReadTimeout,
		WriteTimeout: cfg.Api.WriteTimeout,
		IdleTimeout:  cfg.Api.IdleTimeout,
		ErrorLog:     zap.NewStdLog(logger.Desugar()),
	}

	apiErrors := make(chan error, 1)

	go func() {
		logger.Infow("Starting service", "host", cfg.Api.Host)
		defer logger.Infow("Service stopped", "host", cfg.Api.Host)
		apiErrors <- apiSrv.ListenAndServe()
	}()

	select {
	case err = <-apiErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		logger.Infow("Starting shutdown", "signal", sig)
		defer logger.Infow("Shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), cfg.Api.ShutdownTimeout)
		defer cancel()

		if err = apiSrv.Shutdown(ctx); err != nil {
			_ = apiSrv.Close()
			return fmt.Errorf("cannot shutdown server gracefully: %w", err)
		}
	}

	return nil
}
