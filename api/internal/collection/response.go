package collection

import "encoding/json"

type CollectionResponse struct {
	Slug      string          `json:"slug"`
	Name      string          `json:"name"`
	Fields    json.RawMessage `json:"fields"`
	CreatedAt string          `json:"created_at"`
	UpdatedAt string          `json:"updated_at"`
}

func NewCollectionResponse(c *Collection) *CollectionResponse {
	return &CollectionResponse{
		Slug:      c.Slug,
		Name:      c.Name,
		Fields:    c.Fields,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}
