package web

import (
	"context"
	"encoding/json"
	"net/http"
)

func Response(ctx context.Context, w http.ResponseWriter, statusCode int, data interface{}) error {

	// If there is nothing to marshal then set status code and return.
	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}

	// Marshal data to JSON.
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Ensure content type has been set properly while we know marshaling has succeeded.
	w.Header().Set("Content-Type", "application/json")

	// Write the status code to the response.
	w.WriteHeader(statusCode)

	// Send result back to the client.
	_, err = w.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}
