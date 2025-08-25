package schema_generator

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/lib/pq"
	"github.com/mimsy-cms/mimsy/pkg/mimsy_schema"
)

type DatabaseCollection struct {
	collection *mimsy_schema.Collection
}

func GetPrefixedTableName(relationName string) (string, error) {
	if regexp.MustCompile(`<builtins\.[a-zA-Z0-9]+>`).MatchString(relationName) {
		switch relationName {
		case "<builtins.user>":
			return "mimsy_internal.\"user\"", nil
		case "<builtins.media>":
			return "mimsy_internal.\"media\"", nil
		default:
			return "", errors.New("unknown builtin reference")
		}
	} else {
		return fmt.Sprintf("mimsy_collections.%s", pq.QuoteIdentifier(relationName)), nil
	}
}

func GetSimpleTableName(relationName string) (string, error) {
	if regexp.MustCompile(`<builtins\.[a-zA-Z0-9]+>`).MatchString(relationName) {
		switch relationName {
		case "<builtins.user>":
			return "user", nil
		case "<builtins.media>":
			return "media", nil
		default:
			return "", errors.New("unknown builtin reference")
		}
	} else {
		return relationName, nil
	}
}

// IsBuiltin returns whether the collection is a builtin mimsy collection or not.
func IsBuiltin(name string) bool {
	return regexp.MustCompile(`<builtins\.[a-zA-Z0-9]+>`).MatchString(name)
}

func RemoveSchemaFromReference(relationName string) string {
	re := regexp.MustCompile(`(?m)[A-Za-z_\-0-9]+\.\"?([a-zA-Z_\-0-9]+)\"?`)

	return re.ReplaceAllString(relationName, "$1")
}

func (dc *DatabaseCollection) GetCollectionTableName() string {
	return dc.collection.Name
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

	relatesTo, err := GetSimpleTableName(c.RelatesTo)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s_%s_relation_%s", collectionTableName, fieldName, relatesTo), nil
}
