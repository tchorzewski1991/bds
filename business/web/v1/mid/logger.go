package mid

import (
	"context"
	"net/http"
	"time"

	"github.com/tchorzewski1991/bds/base/web"
	"go.uber.org/zap"
)

func Logger(logger *zap.SugaredLogger) web.Middleware {

	// This is a web.Middleware func that will be executed
	m := func(handler web.Handler) web.Handler {

		// This is a web.Handler func that will wrap original handler with logging functionality
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			// Try to extract ctx values from ctx.
			v, err := web.GetCtxValues(ctx)
			if err != nil {
				return web.NewShutdownError("cannot fetch values from context")
			}

			logger.Infow("request started",
				"trace_id", v.TraceID,
				"method", r.Method,
				"path", r.URL.Path,
			)

			// Call the next handler in the chain.
			err = handler(ctx, w, r)

			logger.Infow("request ended",
				"trace_id", v.TraceID,
				"method", r.Method,
				"path", r.URL.Path,
				"status", v.StatusCode,
				"time", time.Since(v.Now).String(),
			)

			return err
		}

		return h
	}
	return m
}
