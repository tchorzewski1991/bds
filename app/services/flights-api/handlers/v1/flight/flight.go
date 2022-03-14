package flight

import (
	"context"
	"encoding/json"
	"net/http"
)

// Notes on HTTP handlers:
// - Handlers are presentation layer.
//   They take external input, process it and send the response back to external output.
// - There is a bunch of details we want to keep consistent between each of these handlers
//   like: logging, error handling or JSON marshaling protocol.

func List(_ context.Context, w http.ResponseWriter, _ *http.Request) error {
	type flight struct {
		Identifier string `json:"identifier"`
	}
	list := []flight{
		{
			Identifier: "LH-1111-20220101-GDN-WAW",
		},
	}
	err := json.NewEncoder(w).Encode(list)
	if err != nil {
		return err
	}
	return nil
}
