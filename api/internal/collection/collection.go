package collection

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func DefinitionHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Handler logic for getting collection definition
		slug := r.PathValue("collectionSlug")
		if slug == "" {
			http.Error(w, "Missing slug", http.StatusBadRequest)
			return
		}

		var name string
		var fields json.RawMessage
		var createdAt, updatedAt string
		var createdBy, updatedBy string

		err := db.QueryRow(`SELECT name, fields, created_at, created_by, updated_at, updated_by FROM "collection" WHERE slug = $1`, slug).Scan(&name, &fields, &createdAt, &createdBy, &updatedAt, &updatedBy)
		if err == sql.ErrNoRows {
			http.Error(w, "Collection not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		resp := map[string]interface{}{
			"slug":       slug,
			"name":       name,
			"fields":     json.RawMessage(fields),
			"created_at": createdAt,
			"created_by": createdBy,
			"updated_at": updatedAt,
			"updated_by": updatedBy,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func ItemsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Handler logic for getting collection items
		slug := r.PathValue("collectionSlug")
		if slug == "" {
			http.Error(w, "Missing slug", http.StatusBadRequest)
			return
		}

		var exists bool
		err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM "collection" WHERE slug = $1)`, slug).Scan(&exists)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if !exists {
			http.Error(w, "Collection not found", http.StatusNotFound)
			return
		}

		rows, err := db.Query(`SELECT id, data, slug FROM "collection_item" WHERE collection_slug = $1`, slug)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		type Item struct {
			ID           int             `json:"id"`
			ResourceSlug string          `json:"slug"`
			Data         json.RawMessage `json:"data"`
		}

		var items []Item
		for rows.Next() {
			var item Item
			if err := rows.Scan(&item.ID, &item.Data, &item.ResourceSlug); err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			items = append(items, item)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	}
}
