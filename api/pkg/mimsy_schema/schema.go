package mimsy_schema

import (
	"time"
)

type CollectionFields map[string]SchemaElement

type MimsyConfig struct {
	SchemaPath string `json:"manifestPath"`
	BasePath   string `json:"basePath"`
}

type Schema struct {
	Collections []Collection `json:"collections"`
	GeneratedAt time.Time    `json:"generatedAt"`
}

type Collection struct {
	Name     string           `json:"name"`
	Schema   CollectionFields `json:"schema"`
	IsGlobal bool             `json:"isGlobal,omitempty"`
}

type SchemaElement struct {
	Type      string                `json:"type"`
	RelatesTo string                `json:"relatesTo,omitempty"`
	Options   *SchemaElementOptions `json:"options,omitempty"`
}

type SchemaElementOptions struct {
	Description string                    `json:"description,omitempty"`
	Constraints *SchemaElementConstraints `json:"constraints,omitempty"`
}

type SchemaElementConstraints struct {
	Required  bool `json:"required,omitempty"`
	MinLength int  `json:"minLength,omitempty"`
	MaxLength int  `json:"maxLength,omitempty"`
}

// Helper methods for SchemaElement

// IsRelation returns true if the schema element is a relation type
func (se *SchemaElement) IsRelation() bool {
	return se.Type == "relation" || se.Type == "multi_relation"
}

// IsRequired returns true if the schema element has a required constraint
func (se *SchemaElement) IsRequired() bool {
	if se.Options != nil && se.Options.Constraints != nil {
		return se.Options.Constraints.Required
	}
	return false
}

// GetDescription returns the description of the schema element
func (se *SchemaElement) GetDescription() string {
	if se.Options != nil {
		return se.Options.Description
	}
	return ""
}

// Helper methods for Schema

// GetCollection returns a collection by name, or nil if not found
func (s *Schema) GetCollection(name string) *Collection {
	for i := range s.Collections {
		if s.Collections[i].Name == name {
			return &s.Collections[i]
		}
	}
	return nil
}

// Helper methods for Collection

// GetField returns a schema element by field name, or nil if not found
func (c *Collection) GetField(fieldName string) *SchemaElement {
	if element, exists := c.Schema[fieldName]; exists {
		return &element
	}
	return nil
}

// GetRelationFields returns all fields that are relations
func (c *Collection) GetRelationFields() map[string]SchemaElement {
	relations := make(map[string]SchemaElement)
	for name, element := range c.Schema {
		if element.IsRelation() {
			relations[name] = element
		}
	}
	return relations
}

// GetRequiredFields returns all fields that are required
func (c *Collection) GetRequiredFields() map[string]SchemaElement {
	required := make(map[string]SchemaElement)
	for name, element := range c.Schema {
		if element.IsRequired() {
			required[name] = element
		}
	}
	return required
}
