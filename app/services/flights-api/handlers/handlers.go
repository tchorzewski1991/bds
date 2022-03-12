// Package handlers contains all the handler functions and routes
// supported by the web api.
package handlers

import (
	"expvar"
	"github.com/dimfeld/httptreemux/v5"
	"github.com/tchorzewski1991/fds/app/services/flights-api/handlers/debug/checkgrp"
	v1 "github.com/tchorzewski1991/fds/app/services/flights-api/handlers/v1"
	v2 "github.com/tchorzewski1991/fds/app/services/flights-api/handlers/v2"
	"go.uber.org/zap"
	"net/http"
	"net/http/pprof"
	"os"
)

type DebugMuxConfig struct {
	Build  string
	Logger *zap.SugaredLogger
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

	// Register check group endpoints
	cgh := checkgrp.Handlers{
		Build:  cfg.Build,
		Logger: cfg.Logger,
	}
	mux.HandleFunc("/debug/readiness", cgh.Readiness)
	mux.HandleFunc("/debug/liveness", cgh.Liveness)

	return mux
}

type ApiMuxConfig struct {
	Shutdown chan os.Signal
	Logger   *zap.SugaredLogger
}

func ApiMux(cfg ApiMuxConfig) http.Handler {
	mux := httptreemux.NewContextMux()

	// Load the v1 routes.
	v1.Routes(mux, v1.Config{Logger: cfg.Logger})

	// Load the v2 routes.
	v2.Routes(mux, v2.Config{Logger: cfg.Logger})

	return mux
}
