package flight

import (
	"context"
	"github.com/tchorzewski1991/fds/base/web"
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

// private

type flight struct {
	Identifier string `json:"identifier"`
}

var flights = []flight{
	{
		Identifier: "LH-1111-20220101-GDN-WAW",
	},
}
