package db

import (
	"context"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/tchorzewski1991/bds/business/sys/database"
	"go.uber.org/zap"
)

var ErrBookNotFound = errors.New("book not found")

type Store struct {
	db *database.ExtContext
}

func NewStore(db *sqlx.DB, logger *zap.SugaredLogger) Store {
	return Store{db: database.NewExtContext(db).WithLogger(logger)}
}

func (s Store) QueryByID(ctx context.Context, id int) (Book, error) {
	const q = `select * from books where id = :id`

	ext := s.db.
		WithErrorMapper(database.NewErrorMapper()).
		WithMetric(database.NewHistogram("books", "QueryByID"))

	rows, err := sqlx.NamedQueryContext(ctx, ext, q, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return Book{}, err
	}
	defer rows.Close()

	if !rows.Next() {
		return Book{}, ErrBookNotFound
	}

	var book Book
	err = rows.StructScan(&book)
	if err != nil {
		return Book{}, err
	}

	return book, nil
}

func (s Store) Query(ctx context.Context, page int, rowsPerPage int) ([]Book, error) {
	const q = `select * from books order by id offset :offset rows fetch next :rows_per_page rows only`

	// Ensure page is set correctly
	if page < 1 {
		page = 1
	}

	// Ensure rowsPerPage is set correctly
	if rowsPerPage < 1 || rowsPerPage > 20 {
		rowsPerPage = 20
	}

	ext := s.db.
		WithErrorMapper(database.NewErrorMapper()).
		WithMetric(database.NewHistogram("books", "Query"))

	rows, err := sqlx.NamedQueryContext(ctx, ext, q, map[string]interface{}{
		"offset":        (page - 1) * rowsPerPage,
		"rows_per_page": rowsPerPage,
	})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book

	for rows.Next() {
		var book Book
		err = rows.StructScan(&book)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	return books, nil
}
