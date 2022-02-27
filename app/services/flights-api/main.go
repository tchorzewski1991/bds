package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/ardanlabs/conf/v3"

	"github.com/emadolsky/automaxprocs/maxprocs"
	"github.com/tchorzewski1991/fds/base/logger"
	_ "go.uber.org/automaxprocs"
	"go.uber.org/zap"
)

var build = "develop"
var service = "FLIGHTS-API"

func main() {
	// ================================================================================================================
	// Set GOMAXPROCS

	// Set the correct number of threads for the service
	// based on what is available either by the machine or quotas.
	if _, err := maxprocs.Set(); err != nil {
		fmt.Println("maxprocs: %w", err)
	}
	cpu := runtime.GOMAXPROCS(0)

	// ================================================================================================================
	// Construct application logger

	// Set logger fields common to all logs
	fields := []logger.Field{
		{
			Name: "service", Value: service,
		},
		{
			Name: "build", Value: build,
		},
		{
			Name: "cpu", Value: cpu,
		},
	}
	// Create new logger
	l, err := logger.New(fields...)
	if err != nil {
		fmt.Println("Constructing logger error:", err)
		os.Exit(1)
	}
	// Flush logger buffer
	defer l.Sync()

	// Run application
	if err = run(l); err != nil {
		l.Error(err)
		l.Errorf("Running app error: %s", err)
		os.Exit(1)
	}
}

func run(logger *zap.SugaredLogger) error {
	// ================================================================================================================
	// Configuration
	logger.Info("Parsing config")

	cfg := struct {
		conf.Version
	}{
		Version: conf.Version{
			Build: build,
			Desc:  "Current build version",
		},
	}

	const prefix = "FLIGHTS"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config:  %w", err)
	}

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config output: %w", err)
	}
	logger.Infow("Config parsed", "config", out)

	// ================================================================================================================
	// Starting App

	logger.Info("Starting service")
	defer logger.Info("Service stopped")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)
	<-shutdown

	logger.Info("Stopping service")

	return nil
}
