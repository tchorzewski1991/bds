package mid

import (
	"context"
	"fmt"
	"github.com/tchorzewski1991/fds/base/web"
	"github.com/tchorzewski1991/fds/business/sys/auth"
	v1 "github.com/tchorzewski1991/fds/business/web/v1"
	"net/http"
)

func Authorize(permission string) web.Middleware {

	// m is the middleware function to be executed.
	m := func(handler web.Handler) web.Handler {

		// h is the handler that will be attached in the middleware chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			claims, err := auth.GetClaims(ctx)
			if err != nil {
				err = fmt.Errorf("you are not authorized to perform this action, no claims")
				return v1.NewRequestError(err, http.StatusForbidden)
			}

			err = auth.Authorize(claims, permission)
			if err != nil {
				err = fmt.Errorf("you are not authorized to perform this action, permissions missing")
				return v1.NewRequestError(err, http.StatusForbidden)
			}

			return handler(ctx, w, r)
		}

		return h
	}

	return m
}
