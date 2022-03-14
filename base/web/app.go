package web

import (
	"context"
	"github.com/dimfeld/httptreemux/v5"
	"net/http"
	"os"
	"syscall"
)

// Handler represents type responsible for handling http request.
// It extends signature of http.HandlerFunc with support for context.Context.
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

type App struct {
	mux      *httptreemux.ContextMux
	shutdown chan os.Signal
}

func NewApp(shutdown chan os.Signal) *App {
	return &App{
		mux:      httptreemux.NewContextMux(),
		shutdown: shutdown,
	}
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}

func (a *App) Shutdown() {
	a.shutdown <- syscall.SIGTERM
}

func (a *App) Handle(method string, version string, path string, handler Handler) {
	// Prepare the function to execute for each request.
	// This anonymous func wraps Handler with proper error handling.
	h := func(w http.ResponseWriter, r *http.Request) {
		if err := handler(r.Context(), w, r); err != nil {
			a.Shutdown()
			return
		}
	}
	// Extend path with version if necessary.
	if version != "" {
		path = "/" + version + path
	}
	// Register handler func with the requested method and path.
	a.mux.Handle(method, path, h)
}
