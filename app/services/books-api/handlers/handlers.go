// Package handlers contains all the handler functions and routes
// supported by the web api.
package handlers

import (
	"expvar"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tchorzewski1991/bds/app/services/books-api/handlers/debug"
	v1 "github.com/tchorzewski1991/bds/app/services/books-api/handlers/v1"
	v2 "github.com/tchorzewski1991/bds/app/services/books-api/handlers/v2"
	"github.com/tchorzewski1991/bds/base/web"
	"github.com/tchorzewski1991/bds/business/web/v1/mid"
	"go.uber.org/zap"
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

	// Setup liveness/readiness routes.
	ch := debug.CheckHandler{
		Build:  cfg.Build,
		Logger: cfg.Logger,
		DB:     cfg.DB,
	}
	mux.HandleFunc("/debug/readiness", ch.Readiness)
	mux.HandleFunc("/debug/liveness", ch.Liveness)

	return mux
}

type ApiMuxConfig struct {
	Shutdown chan os.Signal
	Logger   *zap.SugaredLogger
	DB       *sqlx.DB
}

func ApiMux(cfg ApiMuxConfig) http.Handler {

	// Initialize new web app with all necessary dependencies.
	app := web.NewApp(
		cfg.Shutdown,
		mid.Metrics(),
		mid.Logger(cfg.Logger),
		mid.Errors(cfg.Logger),
		mid.Panics(),
	)

	// Setup v1 routes.
	v1.Routes(app, v1.Config{Logger: cfg.Logger, DB: cfg.DB})

	// Setup v2 routes.
	v2.Routes(app, v2.Config{Logger: cfg.Logger})

	return app
}
