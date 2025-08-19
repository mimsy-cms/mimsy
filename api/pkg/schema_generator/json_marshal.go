package schema_generator

import (
	"encoding/json"
	"fmt"
)

// ConstraintWrapper wraps a Constraint with type information for JSON marshaling
type ConstraintWrapper struct {
	Type       string          `json:"type"`
	Constraint json.RawMessage `json:"constraint"`
}

// MarshalJSON implements json.Marshaler for Table
func (t *Table) MarshalJSON() ([]byte, error) {
	type Alias Table
	
	// Convert constraints to wrappers
	wrappers := make([]ConstraintWrapper, len(t.Constraints))
	for i, c := range t.Constraints {
		var constraintType string
		switch c.(type) {
		case *UniqueConstraint:
			constraintType = "unique"
		case *PrimaryKeyConstraint:
			constraintType = "primary_key"
		case *CompositePrimaryKeyConstraint:
			constraintType = "composite_primary_key"
		case *ForeignKeyConstraint:
			constraintType = "foreign_key"
		default:
			return nil, fmt.Errorf("unknown constraint type: %T", c)
		}
		
		data, err := json.Marshal(c)
		if err != nil {
			return nil, err
		}
		
		wrappers[i] = ConstraintWrapper{
			Type:       constraintType,
			Constraint: data,
		}
	}
	
	return json.Marshal(&struct {
		*Alias
		Constraints []ConstraintWrapper `json:"Constraints"`
	}{
		Alias:       (*Alias)(t),
		Constraints: wrappers,
	})
}

// UnmarshalJSON implements json.Unmarshaler for Table
func (t *Table) UnmarshalJSON(data []byte) error {
	type Alias Table
	aux := &struct {
		*Alias
		Constraints []ConstraintWrapper `json:"Constraints"`
	}{
		Alias: (*Alias)(t),
	}
	
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	
	// Convert wrappers back to constraints
	t.Constraints = make([]Constraint, len(aux.Constraints))
	for i, wrapper := range aux.Constraints {
		switch wrapper.Type {
		case "unique":
			var c UniqueConstraint
			if err := json.Unmarshal(wrapper.Constraint, &c); err != nil {
				return err
			}
			t.Constraints[i] = &c
		case "primary_key":
			var c PrimaryKeyConstraint
			if err := json.Unmarshal(wrapper.Constraint, &c); err != nil {
				return err
			}
			t.Constraints[i] = &c
		case "composite_primary_key":
			var c CompositePrimaryKeyConstraint
			if err := json.Unmarshal(wrapper.Constraint, &c); err != nil {
				return err
			}
			t.Constraints[i] = &c
		case "foreign_key":
			var c ForeignKeyConstraint
			if err := json.Unmarshal(wrapper.Constraint, &c); err != nil {
				return err
			}
			t.Constraints[i] = &c
		default:
			return fmt.Errorf("unknown constraint type: %s", wrapper.Type)
		}
	}
	
	return nil
}