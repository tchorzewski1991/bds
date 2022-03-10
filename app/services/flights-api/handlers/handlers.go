// Package handlers contains all the handler functions and routes
// supported by the web api.
package handlers

import (
	"expvar"
	"github.com/dimfeld/httptreemux/v5"
	"github.com/tchorzewski1991/fds/app/services/flights-api/handlers/debug/checkgrp"
	"go.uber.org/zap"
	"net/http"
	"net/http/pprof"
	"os"
)

func DebugMux() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/vars", expvar.Handler())

	return mux
}

type ApiMuxConfig struct {
	Shutdown chan os.Signal
	Logger   *zap.SugaredLogger
}

func ApiMux(_ ApiMuxConfig) http.Handler {
	mux := httptreemux.NewContextMux()

	return mux
}
