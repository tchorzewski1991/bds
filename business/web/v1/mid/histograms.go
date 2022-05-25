package mid

import (
	"context"
	"github.com/tchorzewski1991/bds/base/web"
	"github.com/tchorzewski1991/bds/business/sys/metrics"
	"net/http"
)

func Histograms() web.Middleware {

	// m is the middleware function to be executed.
	m := func(handler web.Handler) web.Handler {

		// h is the handler that will be attached in the middleware chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			// Try to extract ctx values from ctx.
			v, err := web.GetCtxValues(ctx)
			if err != nil {
				return web.NewShutdownError("cannot fetch values from context")
			}

			// Prepare and send http histogram
			m := metrics.HttpHistogram(r, v)
			defer m.Send()

			return handler(ctx, w, r)
		}

		return h
	}

	return m
}
