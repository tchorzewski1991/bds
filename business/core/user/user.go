package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	uid "github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/tchorzewski1991/bds/business/core/user/db"
	"github.com/tchorzewski1991/bds/business/sys/auth"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	ErrNotFound         = errors.New("user not found")
	ErrInvalidUUID      = errors.New("UUID is not valid")
	ErrInvalidEmail     = errors.New("email is not valid")
	ErrNotAuthenticated = errors.New("user not authenticated")
)

// Core manages the set of APIs for user access.
// Notes:
// Core does not maintain any state, we should use value semantic.
// Core is responsible for validating user data.
// Core is responsible for persisting user data.
type Core struct {
	store db.Store
}

// NewCore constructs a Core for user api access.
func NewCore(sqlDB *sqlx.DB, logger *zap.SugaredLogger) Core {
	return Core{store: db.NewStore(sqlDB, logger)}
}

func (c Core) QueryByUUID(ctx context.Context, uuid string) (User, error) {
	err := checkUUID(uuid)
	if err != nil {
		return User{}, ErrInvalidUUID
	}

	user, err := c.store.QueryByUUID(ctx, uuid)
	if err != nil {
		if errors.Is(err, db.ErrUserNotFound) {
			return User{}, ErrNotFound
		}
		return User{}, fmt.Errorf("query failed: %w", err)
	}

	return User{
		UUID:        user.UUID,
		Email:       user.Email,
		Permissions: user.Permissions,
	}, nil
}

func (c Core) Authenticate(ctx context.Context, email, pass string) (auth.Claims, error) {
	err := checkEmail(email)
	if err != nil {
		return auth.Claims{}, ErrInvalidEmail
	}

	user, err := c.store.QueryByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, db.ErrUserNotFound) {
			return auth.Claims{}, ErrNotFound
		}
		return auth.Claims{}, fmt.Errorf("authenticate failed: %w", err)
	}

	err = checkPass(user, pass)
	if err != nil {
		return auth.Claims{}, ErrNotAuthenticated
	}

	return auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "bds-api",
			Subject:   user.UUID,
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Permissions: user.Permissions,
	}, nil
}

// private

func checkUUID(uuid string) error {
	_, err := uid.Parse(uuid)
	if err != nil {
		return err
	}
	return nil
}

func checkPass(user db.User, pass string) error {
	err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(pass))
	if err != nil {
		return err
	}
	return nil
}

// TODO
func checkEmail(_ string) error {
	return nil
}
