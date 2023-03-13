package v1

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/tchorzewski1991/bds/base/web"
	"github.com/tchorzewski1991/bds/business/core/user"
	"github.com/tchorzewski1991/bds/business/sys/auth"
	"github.com/tchorzewski1991/bds/business/web/v1"
	"net/http"
)

type userHandler struct {
	user user.Core
}

func (h userHandler) Token(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	email, pass, ok := r.BasicAuth()
	if !ok {
		return v1.NewRequestError(errors.New("user email or password is missing"), http.StatusUnauthorized)
	}

	claims, err := h.user.Authenticate(ctx, email, pass)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrNotFound):
			return v1.NewRequestError(err, http.StatusNotFound)
		case errors.Is(err, user.ErrNotAuthenticated):
			return v1.NewRequestError(err, http.StatusUnauthorized)
		default:
			return fmt.Errorf("authenticate user err: %w", err)
		}
	}

	tkn, err := auth.GenerateToken(claims)
	if err != nil {
		return fmt.Errorf("generate token err: %w", err)
	}

	return web.Response(ctx, w, http.StatusOK, struct {
		Token string `json:"token"`
	}{tkn})
}

func (h userHandler) Profile(ctx context.Context, w http.ResponseWriter, _ *http.Request) error {

	// Get claims out of the ctx. At this point we should always have them available
	// as user has already been authenticated and authorized.
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1.NewRequestError(err, http.StatusForbidden)
	}

	usr, err := h.user.QueryByUUID(ctx, claims.Subject)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return v1.NewRequestError(err, http.StatusNotFound)
		}
		return fmt.Errorf("get user profile err: %w", err)
	}

	return web.Response(ctx, w, http.StatusOK, struct {
		UUID  string `json:"uuid"`
		Email string `json:"email"`
	}{
		UUID:  usr.UUID,
		Email: usr.Email,
	})
}
