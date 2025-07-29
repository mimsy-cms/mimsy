package collection

import (
	"database/sql"
	"net/http"
)

func DefinitionHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Handler logic for getting collection definition
	}
}

func ItemsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Handler logic for getting collection items
	}
}
