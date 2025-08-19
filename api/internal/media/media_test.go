package media_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/mimsy-cms/mimsy/internal/auth"
	"github.com/mimsy-cms/mimsy/internal/media"
	mocks "github.com/mimsy-cms/mimsy/internal/mocks/media"
	storageMocks "github.com/mimsy-cms/mimsy/internal/mocks/storage"
)

// =================================================================================================
// Helper Functions
// =================================================================================================

func createMultipartRequest(t *testing.T, url, fieldName, fileName, content string) *http.Request {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreatePart(textproto.MIMEHeader{
		"Content-Disposition": []string{`form-data; name="` + fieldName + `"; filename="` + fileName + `"`},
		"Content-Type":        []string{"image/jpeg"},
	})
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}

	_, err = part.Write([]byte(content))
	if err != nil {
		t.Fatalf("failed to write to form file: %v", err)
	}

	writer.Close()

	req := httptest.NewRequest("POST", url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req
}

func addUserToContext(req *http.Request, user *auth.User) *http.Request {
	ctx := context.WithValue(req.Context(), auth.UserContextKey, user)
	return req.WithContext(ctx)
}

func executeRequest(handler http.Handler, req *http.Request, t *testing.T) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	return w
}

func createMockUser() *auth.User {
	return &auth.User{
		ID:      1,
		Email:   "admin@example.com",
		IsAdmin: true,
	}
}

func createMockMedia() *media.Media {
	id := uuid.New()
	return &media.Media{
		Id:           1,
		Uuid:         id,
		Name:         "test-image.jpg",
		ContentType:  "image/jpeg",
		CreatedAt:    time.Now(),
		Size:         1024,
		UploadedById: 1,
	}
}

func createMockCreateParams() *media.CreateMediaParams {
	id := uuid.New()
	return &media.CreateMediaParams{
		Uuid:         id,
		Name:         "test-image.jpg",
		ContentType:  "image/jpeg",
		Size:         1024,
		UploadedById: 1,
	}
}

// =================================================================================================
// Handler Tests - Upload
// =================================================================================================

func TestUpload_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockMediaService(ctrl)
	handler := media.NewHandler(mockService)

	user := createMockUser()
	mockMedia := createMockMedia()

	mockService.EXPECT().
		Upload(gomock.Any(), gomock.Any(), "image/jpeg", user).
		Return(mockMedia, nil)

	req := createMultipartRequest(t, "/media/upload", "file", "test-image.jpg", "fake image content")
	req = addUserToContext(req, user)

	w := executeRequest(http.HandlerFunc(handler.Upload), req, t)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status Created, got %v", w.Code)
	}
}

func TestUpload_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockMediaService(ctrl)
	handler := media.NewHandler(mockService)

	req := createMultipartRequest(t, "/media/upload", "file", "test-image.jpg", "fake image content")
	// No user in context

	w := executeRequest(http.HandlerFunc(handler.Upload), req, t)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status Unauthorized, got %v", w.Code)
	}
}

func TestUpload_NoFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockMediaService(ctrl)
	handler := media.NewHandler(mockService)

	user := createMockUser()

	// Create request without file
	req := httptest.NewRequest("POST", "/media/upload", strings.NewReader("no file"))
	req.Header.Set("Content-Type", "application/json")
	req = addUserToContext(req, user)

	w := executeRequest(http.HandlerFunc(handler.Upload), req, t)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status Bad Request, got %v", w.Code)
	}

	if !strings.Contains(w.Body.String(), "Failed to get file from form") {
		t.Errorf("expected error message about file, got: %s", w.Body.String())
	}
}

func TestUpload_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockMediaService(ctrl)
	handler := media.NewHandler(mockService)

	user := createMockUser()

	mockService.EXPECT().
		Upload(gomock.Any(), gomock.Any(), "image/jpeg", user).
		Return(nil, errors.New("upload failed"))

	req := createMultipartRequest(t, "/media/upload", "file", "test-image.jpg", "fake image content")
	req = addUserToContext(req, user)

	w := executeRequest(http.HandlerFunc(handler.Upload), req, t)

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

	mockService := mocks.NewMockMediaService(ctrl)
	handler := media.NewHandler(mockService)

	mockMediaList := []media.Media{*createMockMedia()}

	mockService.EXPECT().
		FindAll(gomock.Any()).
		Return(mockMediaList, nil)

	mockService.EXPECT().
		GetTemporaryURL(gomock.Any(), &mockMediaList[0]).
		Return("https://example.com/temp-url", nil)

	req := httptest.NewRequest("GET", "/media", nil)

	w := executeRequest(http.HandlerFunc(handler.FindAll), req, t)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status OK, got %v", w.Code)
	}

	var response []media.MediaResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(response) != 1 {
		t.Errorf("expected 1 media item, got %v", len(response))
	}

	if response[0].URL != "https://example.com/temp-url" {
		t.Errorf("expected URL to be set, got %v", response[0].URL)
	}
}

func TestFindAll_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockMediaService(ctrl)
	handler := media.NewHandler(mockService)

	mockService.EXPECT().
		FindAll(gomock.Any()).
		Return(nil, errors.New("database error"))

	req := httptest.NewRequest("GET", "/media", nil)

	w := executeRequest(http.HandlerFunc(handler.FindAll), req, t)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status Internal Server Error, got %v", w.Code)
	}
}

// =================================================================================================
// Handler Tests - GetById
// =================================================================================================

func TestGetById_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockMediaService(ctrl)
	handler := media.NewHandler(mockService)

	mockMedia := createMockMedia()

	mockService.EXPECT().
		GetById(gomock.Any(), int64(1)).
		Return(mockMedia, nil)

	mockService.EXPECT().
		GetTemporaryURL(gomock.Any(), mockMedia).
		Return("https://example.com/temp-url", nil)

	req := httptest.NewRequest("GET", "/media/1", nil)
	req.SetPathValue("id", "1")

	w := executeRequest(http.HandlerFunc(handler.GetById), req, t)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status OK, got %v", w.Code)
	}

	var response media.MediaResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response.Id != 1 {
		t.Errorf("expected media ID 1, got %v", response.Id)
	}

	if response.URL != "https://example.com/temp-url" {
		t.Errorf("expected URL to be set, got %v", response.URL)
	}
}

func TestGetById_InvalidId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockMediaService(ctrl)
	handler := media.NewHandler(mockService)

	req := httptest.NewRequest("GET", "/media/invalid", nil)
	req.SetPathValue("id", "invalid")

	w := executeRequest(http.HandlerFunc(handler.GetById), req, t)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status Bad Request, got %v", w.Code)
	}
}

func TestGetById_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockMediaService(ctrl)
	handler := media.NewHandler(mockService)

	mockService.EXPECT().
		GetById(gomock.Any(), int64(999)).
		Return(nil, errors.New("not found"))

	req := httptest.NewRequest("GET", "/media/999", nil)
	req.SetPathValue("id", "999")

	w := executeRequest(http.HandlerFunc(handler.GetById), req, t)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status Not Found, got %v", w.Code)
	}
}

// =================================================================================================
// Handler Tests - Delete
// =================================================================================================

func TestDelete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockMediaService(ctrl)
	handler := media.NewHandler(mockService)

	mockMedia := createMockMedia()

	mockService.EXPECT().
		GetById(gomock.Any(), int64(1)).
		Return(mockMedia, nil)

	mockService.EXPECT().
		Delete(gomock.Any(), mockMedia).
		Return(nil)

	req := httptest.NewRequest("DELETE", "/media/1", nil)
	req.SetPathValue("id", "1")

	w := executeRequest(http.HandlerFunc(handler.Delete), req, t)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status No Content, got %v", w.Code)
	}
}

func TestDelete_InvalidId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockMediaService(ctrl)
	handler := media.NewHandler(mockService)

	req := httptest.NewRequest("DELETE", "/media/invalid", nil)
	req.SetPathValue("id", "invalid")

	w := executeRequest(http.HandlerFunc(handler.Delete), req, t)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status Bad Request, got %v", w.Code)
	}
}

func TestDelete_MediaNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockMediaService(ctrl)
	handler := media.NewHandler(mockService)

	mockService.EXPECT().
		GetById(gomock.Any(), int64(999)).
		Return(nil, errors.New("not found"))

	req := httptest.NewRequest("DELETE", "/media/999", nil)
	req.SetPathValue("id", "999")

	w := executeRequest(http.HandlerFunc(handler.Delete), req, t)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status Not Found, got %v", w.Code)
	}
}

func TestDelete_MediaReferenced(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockMediaService(ctrl)
	handler := media.NewHandler(mockService)

	mockMedia := createMockMedia()

	mockService.EXPECT().
		GetById(gomock.Any(), int64(1)).
		Return(mockMedia, nil)

	mockService.EXPECT().
		Delete(gomock.Any(), mockMedia).
		Return(media.ErrMediaReferenced)

	req := httptest.NewRequest("DELETE", "/media/1", nil)
	req.SetPathValue("id", "1")

	w := executeRequest(http.HandlerFunc(handler.Delete), req, t)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected status Conflict, got %v", w.Code)
	}
}

func TestDelete_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockMediaService(ctrl)
	handler := media.NewHandler(mockService)

	mockMedia := createMockMedia()

	mockService.EXPECT().
		GetById(gomock.Any(), int64(1)).
		Return(mockMedia, nil)

	mockService.EXPECT().
		Delete(gomock.Any(), mockMedia).
		Return(errors.New("database error"))

	req := httptest.NewRequest("DELETE", "/media/1", nil)
	req.SetPathValue("id", "1")

	w := executeRequest(http.HandlerFunc(handler.Delete), req, t)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status Internal Server Error, got %v", w.Code)
	}
}

// =================================================================================================
// Service Tests
// =================================================================================================

func TestService_Upload_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := storageMocks.NewMockStorage(ctrl)
	mockRepo := mocks.NewMockRepository(ctrl)
	service := media.NewService(mockStorage, mockRepo)

	user := createMockUser()
	mockMedia := createMockMedia()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test-image.jpg")
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}
	part.Write([]byte("fake image content"))
	writer.Close()

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.ParseMultipartForm(32 << 20)
	fileHeader := req.MultipartForm.File["file"][0]

	// Set up mock expectations
	mockRepo.EXPECT().
		FindByName(gomock.Any(), "test-image.jpg").
		Return(nil, nil)

	mockStorage.EXPECT().
		Upload(gomock.Any(), gomock.Any(), gomock.Any(), "image/jpeg").
		Return(nil)

	mockRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(mockMedia, nil)

	result, err := service.Upload(context.Background(), fileHeader, "image/jpeg", user)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Name != "test-image.jpg" {
		t.Errorf("expected name 'test-image.jpg', got %v", result.Name)
	}
}

func TestService_Upload_FileNameConflict(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := storageMocks.NewMockStorage(ctrl)
	mockRepo := mocks.NewMockRepository(ctrl)
	service := media.NewService(mockStorage, mockRepo)

	user := createMockUser()
	existingMedia := createMockMedia()
	newMedia := createMockMedia()
	newMedia.Name = "test-image(1).jpg" // Expected resolved name

	// Create a fake multipart file
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test-image.jpg")
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}
	part.Write([]byte("fake image content"))
	writer.Close()

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.ParseMultipartForm(32 << 20)
	fileHeader := req.MultipartForm.File["file"][0]

	// Mock expectations - first call finds existing file, second call finds none
	mockRepo.EXPECT().
		FindByName(gomock.Any(), "test-image.jpg").
		Return(existingMedia, nil)

	mockRepo.EXPECT().
		FindByName(gomock.Any(), "test-image(1).jpg").
		Return(nil, nil)

	mockStorage.EXPECT().
		Upload(gomock.Any(), gomock.Any(), gomock.Any(), "image/jpeg").
		Return(nil)

	mockRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, params *media.CreateMediaParams) (*media.Media, error) {
			if params.Name != "test-image(1).jpg" {
				t.Errorf("expected resolved name 'test-image(1).jpg', got %v", params.Name)
			}
			return newMedia, nil
		})

	result, err := service.Upload(context.Background(), fileHeader, "image/jpeg", user)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Name != "test-image(1).jpg" {
		t.Errorf("expected resolved name 'test-image(1).jpg', got %v", result.Name)
	}
}

func TestService_Upload_StorageError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := storageMocks.NewMockStorage(ctrl)
	mockRepo := mocks.NewMockRepository(ctrl)
	service := media.NewService(mockStorage, mockRepo)

	user := createMockUser()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test-image.jpg")
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}
	part.Write([]byte("fake image content"))
	writer.Close()

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.ParseMultipartForm(32 << 20)
	fileHeader := req.MultipartForm.File["file"][0]

	mockStorage.EXPECT().
		Upload(gomock.Any(), gomock.Any(), gomock.Any(), "image/jpeg").
		Return(errors.New("storage error"))

	result, err := service.Upload(context.Background(), fileHeader, "image/jpeg", user)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if result != nil {
		t.Errorf("expected nil result on error, got %v", result)
	}
}

func TestService_GetById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := storageMocks.NewMockStorage(ctrl)
	mockRepo := mocks.NewMockRepository(ctrl)
	service := media.NewService(mockStorage, mockRepo)

	mockMedia := createMockMedia()

	mockRepo.EXPECT().
		GetById(gomock.Any(), int64(1)).
		Return(mockMedia, nil)

	result, err := service.GetById(context.Background(), 1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Id != 1 {
		t.Errorf("expected media ID 1, got %v", result.Id)
	}
}

func TestService_FindAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := storageMocks.NewMockStorage(ctrl)
	mockRepo := mocks.NewMockRepository(ctrl)
	service := media.NewService(mockStorage, mockRepo)

	mockMediaList := []media.Media{*createMockMedia()}

	mockRepo.EXPECT().
		FindAll(gomock.Any()).
		Return(mockMediaList, nil)

	result, err := service.FindAll(context.Background())

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != 1 {
		t.Errorf("expected 1 media item, got %v", len(result))
	}
}

func TestService_GetTemporaryURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := storageMocks.NewMockStorage(ctrl)
	mockRepo := mocks.NewMockRepository(ctrl)
	service := media.NewService(mockStorage, mockRepo)

	mockMedia := createMockMedia()
	expectedURL := "https://example.com/temp-url"

	mockStorage.EXPECT().
		GetTemporaryURL(mockMedia.Uuid.String(), gomock.Any()).
		Return(expectedURL, nil)

	result, err := service.GetTemporaryURL(context.Background(), mockMedia)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != expectedURL {
		t.Errorf("expected URL %v, got %v", expectedURL, result)
	}
}

func TestService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := storageMocks.NewMockStorage(ctrl)
	mockRepo := mocks.NewMockRepository(ctrl)
	service := media.NewService(mockStorage, mockRepo)

	mockMedia := createMockMedia()

	mockRepo.EXPECT().
		Delete(gomock.Any(), mockMedia).
		Return(nil)

	mockStorage.EXPECT().
		Delete(gomock.Any(), mockMedia.Uuid.String()).
		Return(nil)

	err := service.Delete(context.Background(), mockMedia)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestService_Delete_MediaReferenced(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := storageMocks.NewMockStorage(ctrl)
	mockRepo := mocks.NewMockRepository(ctrl)
	service := media.NewService(mockStorage, mockRepo)

	mockMedia := createMockMedia()

	mockRepo.EXPECT().
		Delete(gomock.Any(), mockMedia).
		Return(media.ErrMediaReferenced)

	err := service.Delete(context.Background(), mockMedia)

	if err != media.ErrMediaReferenced {
		t.Fatalf("expected ErrMediaReferenced, got %v", err)
	}
}

// =================================================================================================
// Response Tests
// =================================================================================================

func TestNewMediaResponse(t *testing.T) {
	mockMedia := createMockMedia()

	response := media.NewMediaResponse(mockMedia)

	if response.Id != mockMedia.Id {
		t.Errorf("expected ID %v, got %v", mockMedia.Id, response.Id)
	}

	if response.Uuid != mockMedia.Uuid.String() {
		t.Errorf("expected UUID %v, got %v", mockMedia.Uuid.String(), response.Uuid)
	}

	if response.Name != mockMedia.Name {
		t.Errorf("expected name %v, got %v", mockMedia.Name, response.Name)
	}

	if response.ContentType != mockMedia.ContentType {
		t.Errorf("expected content type %v, got %v", mockMedia.ContentType, response.ContentType)
	}

	if response.Size != mockMedia.Size {
		t.Errorf("expected size %v, got %v", mockMedia.Size, response.Size)
	}

	if response.UploadedById != mockMedia.UploadedById {
		t.Errorf("expected uploaded by %v, got %v", mockMedia.UploadedById, response.UploadedById)
	}

	expectedTime := mockMedia.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
	if response.CreatedAt != expectedTime {
		t.Errorf("expected created at %v, got %v", expectedTime, response.CreatedAt)
	}
}
