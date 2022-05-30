// Package handlers contains all the handler functions and routes
// supported by the web api.
package handlers

import (
	"expvar"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tchorzewski1991/bds/app/services/books-api/handlers/debug/checkgrp"
	v1 "github.com/tchorzewski1991/bds/app/services/books-api/handlers/v1"
	v2 "github.com/tchorzewski1991/bds/app/services/books-api/handlers/v2"
	"github.com/tchorzewski1991/bds/base/web"
	"github.com/tchorzewski1991/bds/business/web/v1/mid"
	"go.uber.org/zap"
	"net/http"
	"net/http/pprof"
	"os"
)

type DebugMuxConfig struct {
	Build  string
	Logger *zap.SugaredLogger
	DB     *sqlx.DB
}

func DebugMux(cfg DebugMuxConfig) http.Handler {
	mux := http.NewServeMux()

	// Register all the std library debug endpoints.
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/vars", expvar.Handler())
	mux.Handle("/metrics", promhttp.Handler())

	// Register check group endpoints
	cgh := checkgrp.Handlers{
		Build:  cfg.Build,
		Logger: cfg.Logger,
		DB:     cfg.DB,
	}
	mux.HandleFunc("/debug/readiness", cgh.Readiness)
	mux.HandleFunc("/debug/liveness", cgh.Liveness)

	return mux
}

type ApiMuxConfig struct {
	Shutdown chan os.Signal
	Logger   *zap.SugaredLogger
	DB       *sqlx.DB
}

func ApiMux(cfg ApiMuxConfig) http.Handler {

	// Notes and inconsistencies:
	// 1. Let's assume web.App uses one logger for both API versions. Question - where should we put this
	//    logging middleware? How to handle situation when we want to use different logging for both API versions?
	//    Shouldn't we create a form of "registry" for version specific middleware?

	app := web.NewApp(
		cfg.Shutdown,
		mid.Metrics(),
		mid.Logger(cfg.Logger),
		mid.Errors(cfg.Logger),
		mid.Panics(),
	)

	// Load the v1 routes.
	v1.Routes(app, v1.Config{Logger: cfg.Logger, DB: cfg.DB})

	// Load the v2 routes.
	v2.Routes(app, v2.Config{Logger: cfg.Logger})

	return app
}
