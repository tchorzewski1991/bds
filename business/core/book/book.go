package book

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/tchorzewski1991/bds/business/core/book/db"
	"github.com/tchorzewski1991/bds/business/sys/database"
	"go.uber.org/zap"
)

var (
	ErrNotFound  = errors.New("book is not found")
	ErrNotValid  = errors.New("book is not valid")
	ErrNotUnique = errors.New("book is not unique")
)

// Core manages the set of APIs for book access.
// Notes:
// Core does not maintain any state, we should use value semantic.
// Core is responsible for validating book data.
// Core is responsible for persisting book data.
type Core struct {
	store db.Store
}

// NewCore constructs a Core for book api access.
func NewCore(sqlDB *sqlx.DB, logger *zap.SugaredLogger) Core {
	return Core{store: db.NewStore(sqlDB, logger)}
}

func (c Core) QueryByID(ctx context.Context, ID int) (Book, error) {
	book, err := c.store.QueryByID(ctx, ID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return Book{}, ErrNotFound
		}
		return Book{}, fmt.Errorf("query failed: %w", err)
	}
	return convertToBook(book), nil
}

func (c Core) Query(ctx context.Context, page int, rowsPerPage int) ([]Book, error) {
	books, err := c.store.Query(ctx, page, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	return convertToBooks(books), nil
}

func (c Core) Create(ctx context.Context, nb NewBook) (Book, error) {
	book := db.Book{
		Isbn:            nb.Isbn,
		Title:           nb.Title,
		Author:          database.Str(nb.Author),
		PublicationYear: database.Str(nb.PublicationYear),
		Publisher:       database.Str(nb.Publisher),
	}
	err := sanityCheck(book)
	if err != nil {
		return Book{}, fmt.Errorf("create failed: %w", err)
	}

	id, err := c.store.Create(ctx, book)
	if err != nil {
		if errors.Is(err, database.ErrNotUnique) {
			return Book{}, fmt.Errorf("create failed: %w", ErrNotUnique)
		}
		return Book{}, fmt.Errorf("create failed: %w", err)
	}

	book.ID = id

	return convertToBook(book), nil
}

// private

func convertToBooks(books []db.Book) []Book {
	result := make([]Book, len(books))

	for i := 0; i < len(books); i++ {
		result[i] = convertToBook(books[i])
	}

	return result
}

func convertToBook(book db.Book) Book {
	var author *string
	if book.Author.Valid {
		author = &book.Author.String
	}

	var publicationYear *string
	if book.PublicationYear.Valid {
		publicationYear = &book.PublicationYear.String
	}

	var publisher *string
	if book.Publisher.Valid {
		publisher = &book.Publisher.String
	}

	return Book{
		ID:              book.ID,
		Isbn:            book.Isbn,
		Title:           book.Title,
		Author:          author,
		PublicationYear: publicationYear,
		Publisher:       publisher,
	}
}

func sanityCheck(book db.Book) error {
	if book.Title == "" {
		return FieldError{field: "title", err: "can't be blank"}
	}
	if book.Isbn == "" {
		return FieldError{field: "isbn", err: "can't be blank"}
	}
	return nil
}
