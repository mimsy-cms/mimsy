package collection

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/mimsy-cms/mimsy/internal/util"
)

type DefinitionResponse struct {
	Slug      string          `json:"slug"`
	Name      string          `json:"name"`
	Fields    json.RawMessage `json:"fields"`
	CreatedAt string          `json:"created_at"`
	CreatedBy string          `json:"created_by"`
	UpdatedAt string          `json:"updated_at"`
	UpdatedBy string          `json:"updated_by"`
}

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

		util.JSON(w, http.StatusOK, DefinitionResponse{
			Slug:      slug,
			Name:      name,
			Fields:    fields,
			CreatedAt: createdAt,
			CreatedBy: createdBy,
			UpdatedAt: updatedAt,
			UpdatedBy: updatedBy,
		})
	}
}
