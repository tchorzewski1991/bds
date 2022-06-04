package web

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Decode(r *http.Request, dest interface{}) error {
	err := json.NewDecoder(r.Body).Decode(dest)
	if err != nil {
		return fmt.Errorf("payload not valid: %w", err)
	}
	return nil
}
