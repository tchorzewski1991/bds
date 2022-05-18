package mid

import (
	"context"
	"fmt"
	"github.com/tchorzewski1991/bds/base/web"
	"github.com/tchorzewski1991/bds/business/sys/metrics"
	"net/http"
	"runtime/debug"
)

// Panics recovers from panic() and captures the error, so it is reported in Errors middleware.
func Panics() web.Middleware {

	// m is the middleware function to be executed.
	m := func(handler web.Handler) web.Handler {

		// h is the handler that will be attached in the middleware chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {

			// How it works?
			// Defer runs in between state between function returning and calling function gaining back control.
			// When the panic happens defer + recover is able to overwrite err according to predefined guidelines.
			defer func() {
				if rec := recover(); rec != nil {
					// Stack trace will be provided.
					trace := debug.Stack()
					err = fmt.Errorf("PANIC [%v] TRACE[%s]", rec, string(trace))

					// Increment number of panics.
					metrics.AddPanics(ctx)
				}
			}()

			// Call the next handler. Err variable will be set automatically as it is declared
			// on the method signature level.
			return handler(ctx, w, r)
		}

		return h
	}

	return m
}
