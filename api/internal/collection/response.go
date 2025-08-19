package collection

import "encoding/json"

type CollectionResponse struct {
	Slug           string          `json:"slug"`
	Name           string          `json:"name"`
	Fields         json.RawMessage `json:"fields"`
	CreatedAt      string          `json:"created_at"`
	CreatedBy      int64           `json:"created_by"`
	CreatedByEmail string          `json:"created_by_email"`
	UpdatedAt      string          `json:"updated_at"`
	UpdatedBy      int64           `json:"updated_by,omitempty"`
	UpdatedByEmail string          `json:"updated_by_email,omitempty"`
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
