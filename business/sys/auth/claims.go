package auth

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"strings"
)

var ErrClaimsNotFound = errors.New("claims not found")
var ErrActionNotAllowed = errors.New("action not allowed")

type Claims struct {
	jwt.RegisteredClaims
	Permissions []string
}

type ctxKey int

const key ctxKey = 1

func SetClaims(ctx context.Context, claims Claims) context.Context {
	return context.WithValue(ctx, key, claims)
}

func GetClaims(ctx context.Context) (Claims, error) {
	if v, ok := ctx.Value(key).(Claims); ok {
		return v, nil
	}
	return Claims{}, ErrClaimsNotFound
}

// TODO: Make permission mechanism a bit more flexible

func Authorize(claims Claims, fn func(resource, action string) bool) error {
	for _, permission := range claims.Permissions {
		resource, action := decodePermission(permission)

		if fn(resource, action) {
			return nil
		}
	}
	return ErrActionNotAllowed
}

// private

func decodePermission(permission string) (string, string) {
	parts := strings.Split(permission, ".")

	// Example: 'flights.protected'
	if len(parts) >= 2 {
		resource := parts[0]
		action := parts[1]
		return resource, action
	}

	// Example: 'flights'
	if len(parts) == 1 {
		resource := parts[0]
		return resource, ""
	}

	return "", ""
}
