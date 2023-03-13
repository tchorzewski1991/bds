package db

import (
	"time"

	"github.com/lib/pq"
)

type User struct {
	UUID         string         `db:"uuid"`
	Email        string         `db:"email"`
	Permissions  pq.StringArray `db:"permissions"`
	PasswordHash []byte         `db:"password_hash"`
	CreatedAt    time.Time      `db:"date_created"`
	UpdatedAt    time.Time      `db:"date_updated"`
}
