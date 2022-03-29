package flight

import (
	"context"
	"errors"
	"fmt"
	"github.com/dimfeld/httptreemux/v5"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/tchorzewski1991/fds/base/web"
	"github.com/tchorzewski1991/fds/business/sys/auth"
	v1 "github.com/tchorzewski1991/fds/business/web/v1"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

// Notes on HTTP handlers:
// - Handlers are presentation layer.
//   They take external input, process it and send the response back to external output.
// - There is a bunch of details we want to keep consistent between each of these handlers
//   like: logging, error handling or JSON marshaling protocol.

func List(ctx context.Context, w http.ResponseWriter, _ *http.Request) error {
	err := web.Response(ctx, w, http.StatusOK, flights)
	if err != nil {
		return err
	}

	return nil
}

func QueryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	data := httptreemux.ContextData(r.Context())
	params := data.Params()

	f, err := getFlight(params["id"])
	if err != nil {
		return v1.NewRequestError(err, http.StatusNotFound)
	}

	err = web.Response(ctx, w, http.StatusOK, f)
	if err != nil {
		return v1.NewRequestError(err, http.StatusInternalServerError)
	}

	return nil
}

func Protected(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	// Get claims out of the ctx.
	// At this point we should always have them available.
	// They are set through auth middleware.
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1.NewRequestError(err, http.StatusForbidden)
	}

	// Ensure claims owner is authorized to perform the action on the resource.
	err = auth.Authorize(claims, func(resource, action string) bool {
		return resource == "flights" && action == "protected"
	})
	if err != nil {
		return v1.NewRequestError(err, http.StatusForbidden)
	}

	err = web.Response(ctx, w, http.StatusOK, struct {
		Status string `json:"status"`
	}{"ok"})
	if err != nil {
		return v1.NewRequestError(err, http.StatusInternalServerError)
	}

	return nil
}

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

// private

// Section flight
// TODO: Move to separate package

type flight struct {
	Identifier string `json:"identifier"`
}

var flights = []flight{
	{
		Identifier: "LH-1111-20220101-GDN-WAW",
	},
}

func getFlight(identifier string) (flight, error) {
	for _, f := range flights {
		if f.Identifier == identifier {
			return f, nil
		}
	}
	return flight{}, errors.New("flight not found")
}

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
