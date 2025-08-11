package collection

import "encoding/json"

type CollectionResponse struct {
	Slug      string          `json:"slug"`
	Name      string          `json:"name"`
	Fields    json.RawMessage `json:"fields"`
	CreatedAt string          `json:"created_at"`
	CreatedBy string          `json:"created_by"`
	UpdatedAt string          `json:"updated_at"`
	UpdatedBy *string         `json:"updated_by,omitempty"`
}

func NewCollectionResponse(c *Collection) *CollectionResponse {
	return &CollectionResponse{
		Slug:      c.Slug,
		Name:      c.Name,
		Fields:    c.Fields,
		CreatedAt: c.CreatedAt,
		CreatedBy: c.CreatedBy,
		UpdatedAt: c.UpdatedAt,
		UpdatedBy: c.UpdatedBy,
	}
}
