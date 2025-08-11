package mimsy_schema_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/mimsy-cms/mimsy/pkg/mimsy_schema"
)

func TestSchemaUnmarshal(t *testing.T) {
	jsonData := `{
		"collections": [
			{
				"name": "tags",
				"schema": {
					"name": {
						"type": "string",
						"options": {
							"description": "The name of the tag",
							"constraints": { "minLength": 2, "maxLength": 50 }
						}
					},
					"color": {
						"type": "string",
						"options": {
							"description": "The color of the tag, in **hexadecimal** format",
							"constraints": { "minLength": 6, "maxLength": 6 }
						}
					}
				}
			},
			{
				"name": "posts",
				"schema": {
					"title": {
						"type": "string",
						"options": {
							"description": "The title of the post",
							"constraints": { "minLength": 5, "maxLength": 100 }
						}
					},
					"author": {
						"type": "relation",
						"relatesTo": "User",
						"options": {
							"description": "The author of the post",
							"constraints": { "required": true }
						}
					},
					"tags": {
						"type": "relation",
						"relatesTo": "tags",
						"options": {
							"description": "The tags associated with the post",
							"constraints": { "required": true }
						}
					},
					"coverImage": {
						"type": "relation",
						"relatesTo": "Media",
						"options": {
							"description": "The cover image of the post",
							"constraints": { "required": true }
						}
					}
				}
			}
		],
		"generatedAt": "2025-08-11T07:55:00.170Z"
	}`

	var schema mimsy_schema.Schema
	err := json.Unmarshal([]byte(jsonData), &schema)
	if err != nil {
		t.Fatalf("Failed to unmarshal schema: %v", err)
	}

	// Test basic structure
	if len(schema.Collections) != 2 {
		t.Errorf("Expected 2 collections, got %d", len(schema.Collections))
	}

	// Test generated time
	expectedTime, _ := time.Parse(time.RFC3339, "2025-08-11T07:55:00.170Z")
	if !schema.GeneratedAt.Equal(expectedTime) {
		t.Errorf("Expected generatedAt %v, got %v", expectedTime, schema.GeneratedAt)
	}

	// Test tags collection
	tagsCollection := schema.GetCollection("tags")
	if tagsCollection == nil {
		t.Fatal("Tags collection not found")
	}

	if len(tagsCollection.Schema) != 2 {
		t.Errorf("Expected 2 fields in tags schema, got %d", len(tagsCollection.Schema))
	}

	// Test name field in tags
	nameField := tagsCollection.GetField("name")
	if nameField == nil {
		t.Fatal("Name field not found in tags collection")
	}

	if nameField.Type != "string" {
		t.Errorf("Expected name field type to be 'string', got '%s'", nameField.Type)
	}

	if nameField.GetDescription() != "The name of the tag" {
		t.Errorf("Unexpected description for name field: %s", nameField.GetDescription())
	}

	if nameField.Options.Constraints.MinLength != 2 {
		t.Errorf("Expected minLength 2, got %d", nameField.Options.Constraints.MinLength)
	}

	if nameField.Options.Constraints.MaxLength != 50 {
		t.Errorf("Expected maxLength 50, got %d", nameField.Options.Constraints.MaxLength)
	}

	// Test posts collection
	postsCollection := schema.GetCollection("posts")
	if postsCollection == nil {
		t.Fatal("Posts collection not found")
	}

	if len(postsCollection.Schema) != 4 {
		t.Errorf("Expected 4 fields in posts schema, got %d", len(postsCollection.Schema))
	}

	// Test author field (relation)
	authorField := postsCollection.GetField("author")
	if authorField == nil {
		t.Fatal("Author field not found in posts collection")
	}

	if !authorField.IsRelation() {
		t.Error("Author field should be a relation")
	}

	if authorField.RelatesTo != "User" {
		t.Errorf("Expected author to relate to 'User', got '%s'", authorField.RelatesTo)
	}

	if !authorField.IsRequired() {
		t.Error("Author field should be required")
	}

	// Test relation fields
	relationFields := postsCollection.GetRelationFields()
	if len(relationFields) != 3 {
		t.Errorf("Expected 3 relation fields in posts, got %d", len(relationFields))
	}

	// Test required fields
	requiredFields := postsCollection.GetRequiredFields()
	if len(requiredFields) != 3 {
		t.Errorf("Expected 3 required fields in posts, got %d", len(requiredFields))
	}
}

func TestSchemaMarshal(t *testing.T) {
	schema := mimsy_schema.Schema{
		Collections: []mimsy_schema.Collection{
			{
				Name: "users",
				Schema: map[string]mimsy_schema.SchemaElement{
					"username": {
						Type: "string",
						Options: &mimsy_schema.SchemaElementOptions{
							Description: "The username",
							Constraints: &mimsy_schema.SchemaElementConstraints{
								Required:  true,
								MinLength: 3,
								MaxLength: 20,
							},
						},
					},
					"profile": {
						Type:      "relation",
						RelatesTo: "profiles",
						Options: &mimsy_schema.SchemaElementOptions{
							Description: "User profile",
						},
					},
				},
			},
		},
		GeneratedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	data, err := json.Marshal(schema)
	if err != nil {
		t.Fatalf("Failed to marshal schema: %v", err)
	}

	// Unmarshal back to verify
	var unmarshaled mimsy_schema.Schema
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal marshaled data: %v", err)
	}

	// Verify the data
	if len(unmarshaled.Collections) != 1 {
		t.Errorf("Expected 1 collection, got %d", len(unmarshaled.Collections))
	}

	usersCollection := unmarshaled.GetCollection("users")
	if usersCollection == nil {
		t.Fatal("Users collection not found")
	}

	usernameField := usersCollection.GetField("username")
	if usernameField == nil {
		t.Fatal("Username field not found")
	}

	if !usernameField.IsRequired() {
		t.Error("Username field should be required")
	}

	profileField := usersCollection.GetField("profile")
	if profileField == nil {
		t.Fatal("Profile field not found")
	}

	if !profileField.IsRelation() {
		t.Error("Profile field should be a relation")
	}

	if profileField.RelatesTo != "profiles" {
		t.Errorf("Expected profile to relate to 'profiles', got '%s'", profileField.RelatesTo)
	}
}

func TestEmptySchema(t *testing.T) {
	jsonData := `{
		"collections": [],
		"generatedAt": "2024-01-01T00:00:00Z"
	}`

	var schema mimsy_schema.Schema
	err := json.Unmarshal([]byte(jsonData), &schema)
	if err != nil {
		t.Fatalf("Failed to unmarshal empty schema: %v", err)
	}

	if len(schema.Collections) != 0 {
		t.Errorf("Expected 0 collections, got %d", len(schema.Collections))
	}
}

func TestCollectionWithEmptySchema(t *testing.T) {
	jsonData := `{
		"collections": [
			{
				"name": "empty",
				"schema": {}
			}
		],
		"generatedAt": "2024-01-01T00:00:00Z"
	}`

	var schema mimsy_schema.Schema
	err := json.Unmarshal([]byte(jsonData), &schema)
	if err != nil {
		t.Fatalf("Failed to unmarshal schema with empty collection: %v", err)
	}

	emptyCollection := schema.GetCollection("empty")
	if emptyCollection == nil {
		t.Fatal("Empty collection not found")
	}

	if len(emptyCollection.Schema) != 0 {
		t.Errorf("Expected 0 fields in empty collection, got %d", len(emptyCollection.Schema))
	}
}
