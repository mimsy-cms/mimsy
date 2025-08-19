package collection_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/mimsy-cms/mimsy/internal/collection"
	mocks "github.com/mimsy-cms/mimsy/internal/mocks/collection"
)

// =================================================================================================
// Helper Functions
// =================================================================================================

func newJSONRequest(t *testing.T, method, url, jsonBody string) *http.Request {
	t.Helper()
	req := httptest.NewRequest(method, url, strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func executeRequest(handler http.Handler, req *http.Request, t *testing.T) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	return w
}

func createMockCollection() *collection.Collection {
	fields := map[string]interface{}{
		"title": map[string]interface{}{
			"type":  "text",
			"label": "Title",
		},
		"content": map[string]interface{}{
			"type":  "textarea",
			"label": "Content",
		},
	}

	fieldsJSON, _ := json.Marshal(fields)

	return &collection.Collection{
		Slug:      "test-collection",
		Name:      "Test Collection",
		Fields:    fieldsJSON,
		CreatedAt: "2024-01-01T00:00:00Z",
		CreatedBy: "admin@example.com",
		UpdatedAt: "2024-01-01T00:00:00Z",
		UpdatedBy: nil,
		IsGlobal:  false,
	}
}

func createMockResource() *collection.Resource {
	return &collection.Resource{
		Id:             1,
		Slug:           "test-resource",
		CreatedAt:      time.Now(),
		CreatedBy:      1,
		CreatedByEmail: "admin@example.com",
		UpdatedAt:      time.Now(),
		UpdatedBy:      1,
		UpdatedByEmail: "admin@example.com",
		Collection:     "test-collection",
		Fields: map[string]any{
			"title":   "Test Title",
			"content": "Test Content",
		},
	}
}

// =================================================================================================
// Handler Tests - Definition
// =================================================================================================

func TestDefinition_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := collection.NewHandler(mockService)

	mockCollection := createMockCollection()

	mockService.EXPECT().
		FindBySlug(gomock.Any(), "test-collection").
		Return(mockCollection, nil)

	req := httptest.NewRequest("GET", "/collections/test-collection", nil)
	req.SetPathValue("slug", "test-collection")

	w := executeRequest(http.HandlerFunc(handler.Definition), req, t)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status OK, got %v", w.Code)
	}

	var response collection.CollectionResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response.Slug != "test-collection" {
		t.Errorf("expected slug 'test-collection', got %v", response.Slug)
	}
}

func TestDefinition_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := collection.NewHandler(mockService)

	mockService.EXPECT().
		FindBySlug(gomock.Any(), "nonexistent").
		Return(nil, collection.ErrNotFound)

	req := httptest.NewRequest("GET", "/collections/nonexistent", nil)
	req.SetPathValue("slug", "nonexistent")

	w := executeRequest(http.HandlerFunc(handler.Definition), req, t)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status Not Found, got %v", w.Code)
	}
}

func TestDefinition_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := collection.NewHandler(mockService)

	mockService.EXPECT().
		FindBySlug(gomock.Any(), "test-collection").
		Return(nil, errors.New("database error"))

	req := httptest.NewRequest("GET", "/collections/test-collection", nil)
	req.SetPathValue("slug", "test-collection")

	w := executeRequest(http.HandlerFunc(handler.Definition), req, t)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status Internal Server Error, got %v", w.Code)
	}
}

// =================================================================================================
// Handler Tests - GetResources
// =================================================================================================

func TestGetResources_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := collection.NewHandler(mockService)

	mockCollection := createMockCollection()
	mockResources := []collection.Resource{*createMockResource()}

	mockService.EXPECT().
		FindBySlug(gomock.Any(), "test-collection").
		Return(mockCollection, nil)

	mockService.EXPECT().
		FindResources(gomock.Any(), mockCollection).
		Return(mockResources, nil)

	req := httptest.NewRequest("GET", "/collections/test-collection/resources", nil)
	req.SetPathValue("slug", "test-collection")

	w := executeRequest(http.HandlerFunc(handler.GetResources), req, t)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status OK, got %v", w.Code)
	}

	var response []collection.Resource
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(response) != 1 {
		t.Errorf("expected 1 resource, got %v", len(response))
	}
}

func TestGetResources_CollectionNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := collection.NewHandler(mockService)

	mockService.EXPECT().
		FindBySlug(gomock.Any(), "nonexistent").
		Return(nil, collection.ErrNotFound)

	req := httptest.NewRequest("GET", "/collections/nonexistent/resources", nil)
	req.SetPathValue("slug", "nonexistent")

	w := executeRequest(http.HandlerFunc(handler.GetResources), req, t)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status Not Found, got %v", w.Code)
	}
}

func TestGetResources_FindResourcesError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := collection.NewHandler(mockService)

	mockCollection := createMockCollection()

	mockService.EXPECT().
		FindBySlug(gomock.Any(), "test-collection").
		Return(mockCollection, nil)

	mockService.EXPECT().
		FindResources(gomock.Any(), mockCollection).
		Return(nil, errors.New("database error"))

	req := httptest.NewRequest("GET", "/collections/test-collection/resources", nil)
	req.SetPathValue("slug", "test-collection")

	w := executeRequest(http.HandlerFunc(handler.GetResources), req, t)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status Internal Server Error, got %v", w.Code)
	}
}

// =================================================================================================
// Handler Tests - GetResource
// =================================================================================================

func TestGetResource_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := collection.NewHandler(mockService)

	mockCollection := createMockCollection()
	mockResource := createMockResource()

	mockService.EXPECT().
		FindBySlug(gomock.Any(), "test-collection").
		Return(mockCollection, nil)

	mockService.EXPECT().
		FindResource(gomock.Any(), mockCollection, "test-resource").
		Return(mockResource, nil)

	mockService.EXPECT().
		FindUserEmail(gomock.Any(), mockResource.CreatedBy).
		Return("admin@example.com", nil)

	mockService.EXPECT().
		FindUserEmail(gomock.Any(), mockResource.UpdatedBy).
		Return("admin@example.com", nil)

	req := httptest.NewRequest("GET", "/collections/test-collection/resources/test-resource", nil)
	req.SetPathValue("slug", "test-collection")
	req.SetPathValue("resourceSlug", "test-resource")

	w := executeRequest(http.HandlerFunc(handler.GetResource), req, t)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status OK, got %v", w.Code)
	}

	var response collection.Resource
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response.Slug != "test-resource" {
		t.Errorf("expected slug 'test-resource', got %v", response.Slug)
	}
}

func TestGetResource_ResourceNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := collection.NewHandler(mockService)

	mockCollection := createMockCollection()

	mockService.EXPECT().
		FindBySlug(gomock.Any(), "test-collection").
		Return(mockCollection, nil)

	mockService.EXPECT().
		FindResource(gomock.Any(), mockCollection, "nonexistent").
		Return(nil, collection.ErrNotFound)

	req := httptest.NewRequest("GET", "/collections/test-collection/resources/nonexistent", nil)
	req.SetPathValue("slug", "test-collection")
	req.SetPathValue("resourceSlug", "nonexistent")

	w := executeRequest(http.HandlerFunc(handler.GetResource), req, t)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status Not Found, got %v", w.Code)
	}
}

// =================================================================================================
// Handler Tests - UpdateResource
// =================================================================================================

func TestUpdateResource_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := collection.NewHandler(mockService)

	mockCollection := createMockCollection()
	mockResource := createMockResource()

	contentData := map[string]any{
		"title":   "Updated Title",
		"content": "Updated Content",
	}

	mockService.EXPECT().
		FindBySlug(gomock.Any(), "test-collection").
		Return(mockCollection, nil)

	mockService.EXPECT().
		UpdateResourceContent(gomock.Any(), mockCollection, "test-resource", contentData).
		Return(mockResource, nil)

	reqBody := `{"title":"Updated Title","content":"Updated Content"}`
	req := newJSONRequest(t, "PUT", "/collections/test-collection/resources/test-resource", reqBody)
	req.SetPathValue("slug", "test-collection")
	req.SetPathValue("resourceSlug", "test-resource")

	w := executeRequest(http.HandlerFunc(handler.UpdateResource), req, t)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status OK, got %v", w.Code)
	}
}

func TestUpdateResource_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := collection.NewHandler(mockService)

	req := newJSONRequest(t, "PUT", "/collections/test-collection/resources/test-resource", `{"title":"Updated Title"`) // malformed JSON
	req.SetPathValue("slug", "test-collection")
	req.SetPathValue("resourceSlug", "test-resource")

	w := executeRequest(http.HandlerFunc(handler.UpdateResource), req, t)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status Bad Request, got %v", w.Code)
	}
}

func TestUpdateResource_CollectionNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := collection.NewHandler(mockService)

	mockService.EXPECT().
		FindBySlug(gomock.Any(), "nonexistent").
		Return(nil, collection.ErrNotFound)

	reqBody := `{"title":"Updated Title"}`
	req := newJSONRequest(t, "PUT", "/collections/nonexistent/resources/test-resource", reqBody)
	req.SetPathValue("slug", "nonexistent")
	req.SetPathValue("resourceSlug", "test-resource")

	w := executeRequest(http.HandlerFunc(handler.UpdateResource), req, t)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status Not Found, got %v", w.Code)
	}
}

func TestUpdateResource_ResourceNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := collection.NewHandler(mockService)

	mockCollection := createMockCollection()
	contentData := map[string]any{"title": "Updated Title"}

	mockService.EXPECT().
		FindBySlug(gomock.Any(), "test-collection").
		Return(mockCollection, nil)

	mockService.EXPECT().
		UpdateResourceContent(gomock.Any(), mockCollection, "nonexistent", contentData).
		Return(nil, collection.ErrNotFound)

	reqBody := `{"title":"Updated Title"}`
	req := newJSONRequest(t, "PUT", "/collections/test-collection/resources/nonexistent", reqBody)
	req.SetPathValue("slug", "test-collection")
	req.SetPathValue("resourceSlug", "nonexistent")

	w := executeRequest(http.HandlerFunc(handler.UpdateResource), req, t)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status Not Found, got %v", w.Code)
	}
}

// =================================================================================================
// Handler Tests - DeleteResource
// =================================================================================================

func TestDeleteResource_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := collection.NewHandler(mockService)

	mockCollection := createMockCollection()
	mockResource := createMockResource()

	mockService.EXPECT().
		FindBySlug(gomock.Any(), "test-collection").
		Return(mockCollection, nil)

	mockService.EXPECT().
		FindResource(gomock.Any(), mockCollection, "test-resource").
		Return(mockResource, nil)

	mockService.EXPECT().
		FindUserEmail(gomock.Any(), mockResource.CreatedBy).
		Return("admin@example.com", nil)

	mockService.EXPECT().
		FindUserEmail(gomock.Any(), mockResource.UpdatedBy).
		Return("admin@example.com", nil)

	mockService.EXPECT().
		DeleteResource(gomock.Any(), mockResource).
		Return(nil)

	req := httptest.NewRequest("DELETE", "/collections/test-collection/resources/test-resource", nil)
	req.SetPathValue("slug", "test-collection")
	req.SetPathValue("resourceSlug", "test-resource")

	w := executeRequest(http.HandlerFunc(handler.DeleteResource), req, t)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status No Content, got %v", w.Code)
	}
}

func TestDeleteResource_ResourceNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := collection.NewHandler(mockService)

	mockCollection := createMockCollection()

	mockService.EXPECT().
		FindBySlug(gomock.Any(), "test-collection").
		Return(mockCollection, nil)

	mockService.EXPECT().
		FindResource(gomock.Any(), mockCollection, "nonexistent").
		Return(nil, collection.ErrNotFound)

	req := httptest.NewRequest("DELETE", "/collections/test-collection/resources/nonexistent", nil)
	req.SetPathValue("slug", "test-collection")
	req.SetPathValue("resourceSlug", "nonexistent")

	w := executeRequest(http.HandlerFunc(handler.DeleteResource), req, t)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status Not Found, got %v", w.Code)
	}
}

func TestDeleteResource_DeleteError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := collection.NewHandler(mockService)

	mockCollection := createMockCollection()
	mockResource := createMockResource()

	mockService.EXPECT().
		FindBySlug(gomock.Any(), "test-collection").
		Return(mockCollection, nil)

	mockService.EXPECT().
		FindResource(gomock.Any(), mockCollection, "test-resource").
		Return(mockResource, nil)

	mockService.EXPECT().
		FindUserEmail(gomock.Any(), mockResource.CreatedBy).
		Return("admin@example.com", nil)

	mockService.EXPECT().
		FindUserEmail(gomock.Any(), mockResource.UpdatedBy).
		Return("admin@example.com", nil)

	mockService.EXPECT().
		DeleteResource(gomock.Any(), mockResource).
		Return(errors.New("database error"))

	req := httptest.NewRequest("DELETE", "/collections/test-collection/resources/test-resource", nil)
	req.SetPathValue("slug", "test-collection")
	req.SetPathValue("resourceSlug", "test-resource")

	w := executeRequest(http.HandlerFunc(handler.DeleteResource), req, t)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status Internal Server Error, got %v", w.Code)
	}
}

// =================================================================================================
// Handler Tests - FindAll
// =================================================================================================

func TestFindAll_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := collection.NewHandler(mockService)

	mockCollections := []collection.Collection{*createMockCollection()}

	mockService.EXPECT().
		FindAll(gomock.Any(), &collection.FindAllParams{Search: ""}).
		Return(mockCollections, nil)

	req := httptest.NewRequest("GET", "/collections", nil)

	w := executeRequest(http.HandlerFunc(handler.FindAll), req, t)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status OK, got %v", w.Code)
	}

	var response []collection.CollectionResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(response) != 1 {
		t.Errorf("expected 1 collection, got %v", len(response))
	}
}

func TestFindAll_WithSearch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := collection.NewHandler(mockService)

	mockCollections := []collection.Collection{*createMockCollection()}

	mockService.EXPECT().
		FindAll(gomock.Any(), &collection.FindAllParams{Search: "test"}).
		Return(mockCollections, nil)

	req := httptest.NewRequest("GET", "/collections?q=test", nil)

	w := executeRequest(http.HandlerFunc(handler.FindAll), req, t)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status OK, got %v", w.Code)
	}
}

func TestFindAll_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := collection.NewHandler(mockService)

	mockService.EXPECT().
		FindAll(gomock.Any(), &collection.FindAllParams{Search: ""}).
		Return(nil, errors.New("database error"))

	req := httptest.NewRequest("GET", "/collections", nil)

	w := executeRequest(http.HandlerFunc(handler.FindAll), req, t)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status Internal Server Error, got %v", w.Code)
	}
}

// =================================================================================================
// Handler Tests - FindAllGlobals
// =================================================================================================

func TestFindAllGlobals_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := collection.NewHandler(mockService)

	globalCollection := createMockCollection()
	globalCollection.IsGlobal = true
	mockGlobals := []collection.Collection{*globalCollection}

	mockService.EXPECT().
		FindAllGlobals(gomock.Any(), &collection.FindAllParams{Search: ""}).
		Return(mockGlobals, nil)

	req := httptest.NewRequest("GET", "/globals", nil)

	w := executeRequest(http.HandlerFunc(handler.FindAllGlobals), req, t)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status OK, got %v", w.Code)
	}

	var response []collection.CollectionResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(response) != 1 {
		t.Errorf("expected 1 global collection, got %v", len(response))
	}
}

// =================================================================================================
// Service Tests
// =================================================================================================

func TestService_FindBySlug(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	service := collection.NewService(mockRepo)

	mockCollection := createMockCollection()

	mockRepo.EXPECT().
		FindBySlug(gomock.Any(), "test-collection").
		Return(mockCollection, nil)

	result, err := service.FindBySlug(context.Background(), "test-collection")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Slug != "test-collection" {
		t.Errorf("expected slug 'test-collection', got %v", result.Slug)
	}
}

func TestService_FindResource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	service := collection.NewService(mockRepo)

	mockCollection := createMockCollection()
	mockResource := createMockResource()

	mockRepo.EXPECT().
		FindResource(gomock.Any(), mockCollection, "test-resource").
		Return(mockResource, nil)

	result, err := service.FindResource(context.Background(), mockCollection, "test-resource")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Slug != "test-resource" {
		t.Errorf("expected slug 'test-resource', got %v", result.Slug)
	}
}

// =================================================================================================
// Resource MarshalJSON Tests
// =================================================================================================

func TestResource_MarshalJSON(t *testing.T) {
	resource := createMockResource()

	// Test with byte slice field (simulating JSON data from DB)
	jsonData := `{"nested":"value"}`
	resource.Fields["json_field"] = []byte(jsonData)

	data, err := json.Marshal(resource)
	if err != nil {
		t.Fatalf("failed to marshal resource: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal result: %v", err)
	}

	// Check that byte slice was converted to JSON object
	if jsonField, ok := result["json_field"].(map[string]interface{}); !ok {
		t.Errorf("expected json_field to be object, got %T", result["json_field"])
	} else if jsonField["nested"] != "value" {
		t.Errorf("expected nested value 'value', got %v", jsonField["nested"])
	}

	// Check other fields are preserved
	if result["slug"] != "test-resource" {
		t.Errorf("expected slug 'test-resource', got %v", result["slug"])
	}
}
