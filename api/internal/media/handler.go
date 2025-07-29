package media

import (
	"net/http"
)

type mediaHandler struct {
	mediaService MediaService
}

func NewHandler(mediaService MediaService) *mediaHandler {
	return &mediaHandler{mediaService: mediaService}
}

func (h *mediaHandler) Upload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(256 * 1024) // 256 MB

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

	_, err = h.mediaService.Upload(r.Context(), header, contentType)
	if err != nil {
		http.Error(w, "Failed to upload file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
