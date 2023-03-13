package web

import (
	"context"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/dimfeld/httptreemux/v5"
	"github.com/google/uuid"
)

// Handler represents type responsible for handling http request.
// It extends signature of http.HandlerFunc with support for context.Context.
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

type App struct {
	mux      *httptreemux.ContextMux
	shutdown chan os.Signal
	mw       []Middleware
}

func NewApp(shutdown chan os.Signal, mw ...Middleware) *App {
	return &App{
		mux:      httptreemux.NewContextMux(),
		shutdown: shutdown,
		mw:       mw,
	}
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}

func (a *App) Shutdown() {
	a.shutdown <- syscall.SIGTERM
}

func (a *App) Handle(method string, version string, path string, handler Handler, mw ...Middleware) {

	// Wrap handler specific middleware.
	handler = wrapMiddleware(mw, handler)

	// Wrap handler with application level middleware.
	handler = wrapMiddleware(a.mw, handler)

	// Prepare the function to execute for each request.
	// This anonymous func wraps Handler with proper error handling.
	h := func(w http.ResponseWriter, r *http.Request) {

		// Pull context out of the *http.Request and extend it with
		// custom data expected by other middleware.
		//
		// Context might be a good place for storing values like:
		// - start time of the request
		// - response code
		// - trace ID
		ctx := r.Context()

		ctx = context.WithValue(ctx, key, &CtxValues{
			TraceID: uuid.Must(uuid.NewRandom()).String(),
			Now:     time.Now().UTC(),
		})

		if err := handler(ctx, w, r); err != nil {
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
