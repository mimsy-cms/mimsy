package util

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// DecodeJSON decodes JSON from the request body into the provided type T.
// It returns an error if the content type is not application/json or if decoding fails.
func DecodeJSON[T any](r *http.Request) (*T, error) {
	var t T

	contentType, _, _ := strings.Cut(r.Header.Get("Content-Type"), ";")
	if contentType != "application/json" {
		return nil, fmt.Errorf("expected application/json content type, got %s", contentType)
	}

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		return nil, err
	}

	return &t, nil
}

// JSON writes the provided type T as JSON to the response writer with the specified status code.
// It sets the Content-Type header to application/json and handles any encoding errors.
// It returns an error if encoding fails.
func JSON[T any](w http.ResponseWriter, status int, t T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(t); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}
