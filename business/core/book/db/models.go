package db

import (
	"database/sql"
	"time"
)

type Book struct {
	ID              int            `db:"id"`
	Isbn            string         `db:"isbn"`
	Title           string         `db:"title"`
	Author          sql.NullString `db:"author"`
	PublicationYear sql.NullString `db:"publication_year"`
	Publisher       sql.NullString `db:"publisher"`
	CreatedAt       time.Time      `db:"created_at"`
	UpdatedAt       time.Time      `db:"updated_at"`
}
