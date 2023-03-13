package v1

import (
	"context"
	"fmt"
	"github.com/dimfeld/httptreemux/v5"
	"github.com/pkg/errors"
	"github.com/tchorzewski1991/bds/base/web"
	"github.com/tchorzewski1991/bds/business/core/book"
	"github.com/tchorzewski1991/bds/business/web/v1"
	"net/http"
	"strconv"
)

type bookHandler struct {
	book book.Core
}

func (h bookHandler) Query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
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

	books, err := h.book.Query(ctx, page, rowsPerPage)
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

func (h bookHandler) QueryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	params := httptreemux.ContextParams(r.Context())

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return v1.NewRequestError(fmt.Errorf("id param is not valid: %w", err), http.StatusBadRequest)
	}

	b, err := h.book.QueryByID(ctx, id)
	if err != nil {
		if errors.Is(err, book.ErrNotFound) {
			return v1.NewRequestError(err, http.StatusNotFound)
		}
		return err
	}

	return web.Response(ctx, w, http.StatusOK, b)
}

func (h bookHandler) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var nb book.NewBook
	err := web.Decode(r, &nb)
	if err != nil {
		return v1.NewRequestError(err, http.StatusBadRequest)
	}

	b, err := h.book.Create(ctx, nb)
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
