package util

import (
	"net/http"

	"github.com/gorilla/schema"
)

var (
	decoder = schema.NewDecoder()
)

func init() {
	decoder.SetAliasTag("query")
}

// QueryString decodes the query parameters from the request into the provided struct.
// It returns a pointer to the struct or an error if decoding fails.
func QueryString[T any](r *http.Request) (*T, error) {
	var result T
	if err := decoder.Decode(&result, r.URL.Query()); err != nil {
		return nil, err
	}
	return &result, nil
}
