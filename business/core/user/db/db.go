package db

import (
	"context"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/tchorzewski1991/bds/business/sys/database"
	"go.uber.org/zap"
)

var ErrUserNotFound = errors.New("user not found")

type Store struct {
	db *database.ExtContext
}

func NewStore(db *sqlx.DB, logger *zap.SugaredLogger) Store {
	return Store{db: database.NewExtContext(db).WithLogger(logger)}
}

func (s Store) QueryByUUID(ctx context.Context, uuid string) (User, error) {
	const q = `select * from users where uuid = :uuid`

	data := struct {
		Uuid string `db:"uuid"`
	}{Uuid: uuid}

	ext := s.db.
		WithErrorMapper(database.NewErrorMapper()).
		WithMetric(database.NewMetric("users", "QueryByUUID"))

	rows, err := sqlx.NamedQueryContext(ctx, ext, q, data)
	if err != nil {
		return User{}, err
	}
	defer rows.Close()

	if !rows.Next() {
		return User{}, ErrUserNotFound
	}

	var user User
	err = rows.StructScan(&user)
	if err != nil {
		return User{}, err
	}

	return user, err
}

func (s Store) QueryByEmail(ctx context.Context, email string) (User, error) {
	const q = `select * from users where email = :email`

	data := struct {
		Email string `db:"email"`
	}{Email: email}

	ext := s.db.
		WithErrorMapper(database.NewErrorMapper()).
		WithMetric(database.NewMetric("users", "QueryByEmail"))

	rows, err := sqlx.NamedQueryContext(ctx, ext, q, data)
	if err != nil {
		return User{}, err
	}
	defer rows.Close()

	if !rows.Next() {
		return User{}, ErrUserNotFound
	}

	var user User
	err = rows.StructScan(&user)
	if err != nil {
		return User{}, err
	}

	return user, err
}
