package mid

import (
	"context"
	"errors"
	"github.com/tchorzewski1991/bds/base/web"
	"github.com/tchorzewski1991/bds/business/sys/auth"
	v1 "github.com/tchorzewski1991/bds/business/web/v1"
	"net/http"
	"strings"
)

func Authenticate() web.Middleware {

	// m is the middleware function to be executed.
	m := func(handler web.Handler) web.Handler {

		// h is the handler that will be attached in the middleware chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				err := errors.New("authorization header is not set")
				return v1.NewRequestError(err, http.StatusUnauthorized)
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) < 2 && strings.ToLower(parts[0]) != "bearer" {
				err := errors.New("authorization header has invalid format. Expected: Bearer TOKEN")
				return v1.NewRequestError(err, http.StatusUnauthorized)
			}

			claims, err := auth.ValidateToken(parts[1])
			if err != nil {
				return v1.NewRequestError(err, http.StatusUnauthorized)
			}

			ctx = auth.SetClaims(ctx, claims)

			return handler(ctx, w, r)
		}

		return h
	}

	return m
}
