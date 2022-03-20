package mid

import (
	"context"
	"github.com/tchorzewski1991/fds/base/web"
	"github.com/tchorzewski1991/fds/business/sys/metrics"
	"net/http"
)

func Metrics() web.Middleware {

	// m is the middleware function to be executed.
	m := func(handler web.Handler) web.Handler {

		// h is the handler that will be attached in the middleware chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			// Extend current ctx with support for metrics gathering.
			ctx = metrics.Set(ctx)

			// Call the next handler.
			err := handler(ctx, w, r)

			// Increment number of requests.
			metrics.AddRequests(ctx)

			// Set number of goroutines.
			metrics.SetGoroutines(ctx)

			// Increment number of errors if necessary.
			if err != nil {
				metrics.AddErrors(ctx)
			}

			// Ensure err is returned, so it can be handler further up the chain.
			return err
		}

		return h
	}

	return m
}
