package media

import (
	"log/slog"
	"net/http"

	"github.com/mimsy-cms/mimsy/internal/auth"
	"github.com/mimsy-cms/mimsy/internal/util"
)

type mediaHandler struct {
	mediaService MediaService
}

func NewHandler(mediaService MediaService) *mediaHandler {
	return &mediaHandler{mediaService: mediaService}
}

func (h *mediaHandler) Upload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(256 * 1024) // 256 MB

	user := auth.UserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		http.Error(w, "Content-Type header is missing", http.StatusBadRequest)
		return
	}

	_, err = h.mediaService.Upload(r.Context(), header, contentType, user)
	if err != nil {
		slog.Error("Failed to upload media", "error", err)
		http.Error(w, "Failed to upload file", http.StatusInternalServerError)
		return
	}

	util.JSON(w, http.StatusCreated, struct{}{})
}

func (h *mediaHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	media, err := h.mediaService.FindAll(r.Context())
	if err != nil {
		slog.Error("Failed to retrieve media", "error", err)
		http.Error(w, "Failed to retrieve media", http.StatusInternalServerError)
		return
	}

	var response []MediaResponse
	for _, m := range media {
		res := NewMediaResponse(&m)
		res.URL, _ = h.mediaService.GetTemporaryURL(r.Context(), &m)
		response = append(response, res)
	}

	util.JSON(w, http.StatusOK, response)
}
