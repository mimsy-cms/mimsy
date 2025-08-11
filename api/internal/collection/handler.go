package collection

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/mimsy-cms/mimsy/internal/util"
)

type Handler struct {
	Service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) Definition(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		http.Error(w, "Missing slug", http.StatusBadRequest)
		return
	}

	def, err := h.Service.GetDefinition(r.Context(), slug)
	if err != nil {
		if err == ErrNotFound {
			http.Error(w, "Collection not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(def)
}

func (h *Handler) GetResources(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		http.Error(w, "Missing slug", http.StatusBadRequest)
		return
	}

	resources, err := h.Service.GetResources(r.Context(), slug)
	if err != nil {
		slog.Error("Failed to get resources", "slug", slug, "error", err)
		if err == ErrNotFound {
			http.Error(w, "Collection not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	if err := util.JSON(w, http.StatusOK, resources); err != nil {
		slog.Error("Failed to encode resources", "slug", slug, "error", err)
	}
}

func (h *Handler) GetResource(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	resourceSlug := r.PathValue("resourceSlug")
	if slug == "" || resourceSlug == "" {
		http.Error(w, "Missing slug or resource slug", http.StatusBadRequest)
		return
	}

	resources, err := h.Service.GetResources(r.Context(), slug)
	if err != nil {
		slog.Error("Failed to get resources", "slug", slug, "error", err)
		if err == ErrNotFound {
			http.Error(w, "Collection not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	for _, resource := range resources {
		if resource["slug"] == resourceSlug {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resource)
			return
		}
	}

	http.Error(w, "Resource not found", http.StatusNotFound)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	collections, err := h.Service.List(r.Context())
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(collections)
}
