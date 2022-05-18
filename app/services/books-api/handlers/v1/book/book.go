package book

import (
	"context"
	"errors"
	"github.com/dimfeld/httptreemux/v5"
	"github.com/tchorzewski1991/bds/base/web"
	v1 "github.com/tchorzewski1991/bds/business/web/v1"
	"net/http"
)

// Notes on HTTP handlers:
// - Handlers are presentation layer.
//   They take external input, process it and send the response back to external output.
// - There is a bunch of details we want to keep consistent between each of these handlers
//   like: logging, error handling or JSON marshaling protocol.

func List(ctx context.Context, w http.ResponseWriter, _ *http.Request) error {
	err := web.Response(ctx, w, http.StatusOK, books)
	if err != nil {
		return err
	}

	return nil
}

func QueryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	data := httptreemux.ContextData(r.Context())
	params := data.Params()

	f, err := getBook(params["id"])
	if err != nil {
		return v1.NewRequestError(err, http.StatusNotFound)
	}

	err = web.Response(ctx, w, http.StatusOK, f)
	if err != nil {
		return v1.NewRequestError(err, http.StatusInternalServerError)
	}

	return nil
}

// private

// Section book
// TODO: Move to separate package

type book struct {
	Identifier string `json:"identifier"`
}

var books = []book{
	{
		Identifier: "1111",
	},
}

func getBook(identifier string) (book, error) {
	for _, b := range books {
		if b.Identifier == identifier {
			return b, nil
		}
	}
	return book{}, errors.New("book not found")
}
