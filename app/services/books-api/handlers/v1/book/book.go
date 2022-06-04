package book

import (
	"context"
	"errors"
	"fmt"
	"github.com/dimfeld/httptreemux/v5"
	"github.com/tchorzewski1991/bds/base/web"
	"github.com/tchorzewski1991/bds/business/core/book"
	v1 "github.com/tchorzewski1991/bds/business/web/v1"
	"net/http"
	"strconv"
)

// Notes on HTTP handlers:
// - Handlers are presentation layer.
//   They take external input, process it and send the response back to external output.
// - There is a bunch of details we want to keep consistent between each of these handlers
//   like: logging, error handling or JSON marshaling protocol.

type Handler struct {
	Book book.Core
}

func (h Handler) Query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()

	var err error
	var page int
	var rowsPerPage int

	if v := query.Get("page"); v != "" {
		page, err = strconv.Atoi(v)
		if err != nil {
			return v1.NewRequestError(fmt.Errorf("page param is not valid: %w", err), http.StatusBadRequest)
		}
	}
	if page < 1 {
		page = 1
	}

	if v := query.Get("rows"); v != "" {
		rowsPerPage, err = strconv.Atoi(v)
		if err != nil {
			return v1.NewRequestError(fmt.Errorf("rows param is not valid: %w", err), http.StatusBadRequest)
		}
	}
	if rowsPerPage < 1 || rowsPerPage > 20 {
		rowsPerPage = 20
	}

	books, err := h.Book.Query(ctx, page, rowsPerPage)
	if err != nil {
		return fmt.Errorf("unable to query books: %w", err)
	}

	return web.Response(ctx, w, http.StatusOK, struct {
		Page  int         `json:"page"`
		Rows  int         `json:"rows"`
		Books []book.Book `json:"books"`
	}{
		Page:  page,
		Rows:  rowsPerPage,
		Books: books,
	})
}

func (h Handler) QueryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	params := httptreemux.ContextParams(r.Context())

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return v1.NewRequestError(fmt.Errorf("id param is not valid: %w", err), http.StatusBadRequest)
	}

	b, err := h.Book.QueryByID(ctx, id)
	if err != nil {
		if errors.Is(err, book.ErrNotFound) {
			return v1.NewRequestError(err, http.StatusNotFound)
		}
		return err
	}

	return web.Response(ctx, w, http.StatusOK, b)
}

func (h Handler) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var nb book.NewBook
	err := web.Decode(r, &nb)
	if err != nil {
		return v1.NewRequestError(err, http.StatusBadRequest)
	}

	b, err := h.Book.Create(ctx, nb)
	if err != nil {
		var fieldErr book.FieldError
		if errors.As(err, &fieldErr) {
			return v1.NewRequestError(fieldErr, http.StatusUnprocessableEntity)
		}
		if errors.Is(err, book.ErrNotUnique) {
			return v1.NewRequestError(err, http.StatusConflict)
		}

		return err
	}

	return web.Response(ctx, w, http.StatusCreated, b)
}
