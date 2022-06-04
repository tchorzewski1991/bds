package mid

import (
	"context"
	"github.com/tchorzewski1991/bds/base/web"
	v1 "github.com/tchorzewski1991/bds/business/web/v1"
	"go.uber.org/zap"
	"net/http"
)

func Errors(logger *zap.SugaredLogger) web.Middleware {

	m := func(handler web.Handler) web.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			// We don't want to continue app execution while ctx values are missing.
			// We want service to be gracefully shutdown.
			v, err := web.GetCtxValues(ctx)
			if err != nil {
				return web.NewShutdownError("cannot fetch values out of context")
			}
			// Execute original handler and act on error according to its type.
			err = handler(ctx, w, r)
			if err != nil {
				// Log the error.
				// TODO: what we want to log here?
				logger.Errorw(err.Error(), "trace_id", v.TraceID)

				// We want to have consistent error response for all the v1 endpoints.
				var er v1.ErrorResponse

				// Build out error response.
				switch {
				case v1.IsFieldError(err):
					fErr := v1.GetFieldError(err)
					er = v1.ErrorResponse{
						Err:    fErr.Error(),
						Status: fErr.Status(),
					}
				case v1.IsRequestError(err):
					rErr := v1.GetRequestError(err)
					er = v1.ErrorResponse{
						Err:    rErr.Error(),
						Status: rErr.Status,
					}
				default:
					er = v1.ErrorResponse{
						Err:    http.StatusText(http.StatusInternalServerError),
						Status: http.StatusInternalServerError,
					}
				}

				// Send error response back to the client.
				err = web.Response(ctx, w, er.Status, er)
				if err != nil {
					return err
				}

				// If we receive the shutdown error we need to return it back to
				// the base handler and shut down service.
				if web.IsShutdownError(err) {
					return err
				}
			}

			return nil
		}

		return h
	}

	return m
}
