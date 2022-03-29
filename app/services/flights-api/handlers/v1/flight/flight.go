package flight

import (
	"context"
	"errors"
	"github.com/dimfeld/httptreemux/v5"
	"github.com/tchorzewski1991/fds/base/web"
	"github.com/tchorzewski1991/fds/business/sys/auth"
	v1 "github.com/tchorzewski1991/fds/business/web/v1"
	"net/http"
)

// Notes on HTTP handlers:
// - Handlers are presentation layer.
//   They take external input, process it and send the response back to external output.
// - There is a bunch of details we want to keep consistent between each of these handlers
//   like: logging, error handling or JSON marshaling protocol.

func List(ctx context.Context, w http.ResponseWriter, _ *http.Request) error {
	err := web.Response(ctx, w, http.StatusOK, flights)
	if err != nil {
		return err
	}

	return nil
}

func QueryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	data := httptreemux.ContextData(r.Context())
	params := data.Params()

	f, err := getFlight(params["id"])
	if err != nil {
		return v1.NewRequestError(err, http.StatusNotFound)
	}

	err = web.Response(ctx, w, http.StatusOK, f)
	if err != nil {
		return v1.NewRequestError(err, http.StatusInternalServerError)
	}

	return nil
}

func Protected(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	// Get claims out of the ctx.
	// At this point we should always have them available.
	// They are set through auth middleware.
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1.NewRequestError(err, http.StatusForbidden)
	}

	// Ensure claims owner is authorized to perform the action on the resource.
	err = auth.Authorize(claims, func(resource, action string) bool {
		return resource == "flights" && action == "protected"
	})
	if err != nil {
		return v1.NewRequestError(err, http.StatusForbidden)
	}

	err = web.Response(ctx, w, http.StatusOK, struct {
		Status string `json:"status"`
	}{"ok"})
	if err != nil {
		return v1.NewRequestError(err, http.StatusInternalServerError)
	}

	return nil
}

// private

// Section flight
// TODO: Move to separate package

type flight struct {
	Identifier string `json:"identifier"`
}

var flights = []flight{
	{
		Identifier: "LH-1111-20220101-GDN-WAW",
	},
}

func getFlight(identifier string) (flight, error) {
	for _, f := range flights {
		if f.Identifier == identifier {
			return f, nil
		}
	}
	return flight{}, errors.New("flight not found")
}
