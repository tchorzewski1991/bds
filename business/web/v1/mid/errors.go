package mid

import (
	"context"
	"github.com/tchorzewski1991/fds/base/web"
	v1 "github.com/tchorzewski1991/fds/business/web/v1"
	"go.uber.org/zap"
	"net/http"
)

func Errors(logger *zap.SugaredLogger) web.Middleware {

	m := func(handler web.Handler) web.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			// Try to extract ctx values from ctx.
			// TODO: Implement request shutdown in case of failure.
			v, err := web.GetCtxValues(ctx)
			if err != nil {
				return err
			}
			// Execute handler and handle error according to its type.
			err = handler(ctx, w, r)
			if err != nil {
				logger.Errorw(err.Error(), "trace_id", v.TraceID)

				var er v1.ErrorResponse

				switch {
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

				err = web.Response(ctx, w, er.Status, er)
				if err != nil {
					return err
				}
			}

			return nil
		}

		return h
	}

	return m
}
