package schema_generator

import (
	"errors"
	"fmt"

	"github.com/mimsy-cms/mimsy/pkg/mimsy_schema"
)

type DatabaseCollection struct {
	collection *mimsy_schema.Collection
}

func (dc *DatabaseCollection) GetCollectionTableName() string {
	return fmt.Sprintf("%s", dc.collection.Name)
}

func (dc *DatabaseCollection) GetFieldRelationTableName(field string) (string, error) {
	fieldDef := dc.collection.GetField(field)
	if fieldDef == nil {
		return "", errors.New("field not found")
	}

	return GetRelationTableName(fieldDef, dc.GetCollectionTableName(), field)
}

func GetRelationTableName(c *mimsy_schema.SchemaElement, collectionTableName string, fieldName string) (string, error) {
	if c.Type != "multi_relation" {
		return "", errors.New("not a multi_relation")
	}

	return fmt.Sprintf("%s_%s_relation_%s", collectionTableName, fieldName, c.RelatesTo), nil
}
