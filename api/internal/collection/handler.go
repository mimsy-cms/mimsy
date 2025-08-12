package collection

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/mimsy-cms/mimsy/internal/util"
)

type Handler struct {
	Service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) Definition(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	collection, err := h.Service.FindById(r.Context(), slug)
	if err != nil {
		slog.Error("Failed to get collection definition", "slug", slug, "error", err)
		if err == ErrNotFound {
			http.Error(w, "Collection not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	util.JSON(w, http.StatusOK, NewCollectionResponse(collection))
}

func (h *Handler) GetResources(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	collection, err := h.Service.FindById(r.Context(), slug)
	if err != nil {
		if err == ErrNotFound {
			http.Error(w, "Collection not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	resources, err := h.Service.FindResources(r.Context(), collection)
	if err != nil {
		slog.Error("Failed to get resources", "slug", slug, "error", err)
		if err == ErrNotFound {
			http.Error(w, "Collection not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	util.JSON(w, http.StatusOK, resources)
}

func (h *Handler) UpdateResource(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	resourceSlug := r.PathValue("resourceSlug")

	var contentData map[string]any
	if err := json.NewDecoder(r.Body).Decode(&contentData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	collection, err := h.Service.FindById(r.Context(), slug)
	if err != nil {
		slog.Error("Failed to get collection", "slug", slug, "error", err)
		if err == ErrNotFound {
			http.Error(w, "Collection not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	updatedResource, err := h.Service.UpdateResourceContent(r.Context(), collection, resourceSlug, contentData)
	if err != nil {
		slog.Error("Failed to update resource", "slug", slug, "resourceSlug", resourceSlug, "error", err)
		if err == ErrNotFound {
			http.Error(w, "Resource not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	util.JSON(w, http.StatusOK, updatedResource)
}

func (h *Handler) GetResource(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	resourceSlug := r.PathValue("resourceSlug")

	collection, err := h.Service.FindById(r.Context(), slug)
	if err != nil {
		slog.Error("Failed to get collection", "slug", slug, "error", err)
		if err == ErrNotFound {
			http.Error(w, "Collection not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	resource, err := h.Service.FindResource(r.Context(), collection, resourceSlug)
	if err != nil {
		slog.Error("Failed to get resource", "slug", slug, "resourceSlug", resourceSlug, "error", err)
		if err == ErrNotFound {
			http.Error(w, "Resource not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	util.JSON(w, http.StatusOK, resource)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	collections, err := h.Service.FindAll(r.Context())
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := make([]CollectionResponse, len(collections))
	for i, collection := range collections {
		response[i] = *NewCollectionResponse(&collection)
	}

	util.JSON(w, http.StatusOK, response)
}

func (h *Handler) ListGlobals(w http.ResponseWriter, r *http.Request) {
	globals, err := h.Service.ListGlobals(r.Context())
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	util.JSON(w, http.StatusOK, globals)
}
