package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/tchorzewski1991/fds/base/web"
	"github.com/tchorzewski1991/fds/business/sys/auth"
	v1 "github.com/tchorzewski1991/fds/business/web/v1"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func Token(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	_, err := web.GetCtxValues(ctx)
	if err != nil {
		return web.NewShutdownError("cannot fetch values out of context")
	}

	name, pass, ok := r.BasicAuth()
	if !ok {
		return v1.NewRequestError(errors.New("user name or password is missing"), http.StatusUnauthorized)
	}

	usr, err := getUser(name, pass)
	if err != nil {
		switch {
		case errors.Is(err, ErrUserNotFound):
			return v1.NewRequestError(err, http.StatusNotFound)
		case errors.Is(err, ErrUserNotAuthenticated):
			return v1.NewRequestError(err, http.StatusUnauthorized)
		default:
			return fmt.Errorf("authentication err: %w", err)
		}
	}

	tkn, err := auth.GenerateToken(auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "fds-api",
			Subject:   usr.uuid,
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Permissions: usr.permissions,
	})
	if err != nil {
		return fmt.Errorf("generating token err: %w", err)
	}

	return web.Response(ctx, w, http.StatusOK, struct {
		Token string `json:"token"`
	}{tkn})
}

// Private

// Section user
// TODO: Move to separate package

var ErrUserNotFound = errors.New("user not found")
var ErrUserNotAuthenticated = errors.New("user not authenticated")

type user struct {
	uuid        string
	name        string
	pass        []byte
	permissions []string
}

var users []user

func init() {
	uid := uuid.NewString()
	pass, _ := bcrypt.GenerateFromPassword([]byte("fds_api_pass"), bcrypt.DefaultCost)
	users = append(users, user{
		uuid:        uid,
		name:        "fds_api_user",
		pass:        pass,
		permissions: []string{"flights.protected"},
	})
}

func getUser(name, pass string) (user, error) {
	for _, u := range users {
		if u.name == name {
			err := bcrypt.CompareHashAndPassword(u.pass, []byte(pass))
			if err != nil {
				return user{}, ErrUserNotAuthenticated
			}
			return u, nil
		}
	}
	return user{}, ErrUserNotFound
}
