package collection

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/mimsy-cms/mimsy/internal/auth"
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

	collection, err := h.Service.FindBySlug(r.Context(), slug)
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

	collection, err := h.Service.FindBySlug(r.Context(), slug)
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
			http.Error(w, "Resources not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	util.JSON(w, http.StatusOK, resources)
}

func (h *Handler) UpdateResource(w http.ResponseWriter, r *http.Request) {
	user := auth.RequestUser(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	slug := r.PathValue("slug")
	resourceSlug := r.PathValue("resourceSlug")

	var contentData map[string]any
	if err := json.NewDecoder(r.Body).Decode(&contentData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	collection, err := h.Service.FindBySlug(r.Context(), slug)
	if err != nil {
		slog.Error("Failed to get collection", "slug", slug, "error", err)
		if err == ErrNotFound {
			http.Error(w, "Collection not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	updatedResource, err := h.Service.UpdateResource(r.Context(), collection, resourceSlug, contentData)
	if err != nil {
		if err == ErrNotFound {
			createdResource, createErr := h.Service.CreateResource(r.Context(), collection, resourceSlug, user.ID, contentData)
			if createErr != nil {
				slog.Error("Failed to create resource", "slug", slug, "resourceSlug", resourceSlug, "error", createErr)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			util.JSON(w, http.StatusCreated, createdResource)
			return
		}
	}

	util.JSON(w, http.StatusOK, updatedResource)
}

type CreateResourceRequest struct {
	Slug string `json:"slug"`
}

func (h *Handler) CreateResource(w http.ResponseWriter, r *http.Request) {
	user := auth.RequestUser(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	collectionSlug := r.PathValue("slug")

	var fields map[string]any
	if err := json.NewDecoder(r.Body).Decode(&fields); err != nil {
		slog.Error("Failed to decode resource content", "error", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	slug, ok := fields["slug"].(string)
	if !ok || slug == "" {
		slog.Error("Failed to get slug from request")
		http.Error(w, "Invalid slug", http.StatusBadRequest)
		return
	}

	collection, err := h.Service.FindBySlug(r.Context(), collectionSlug)
	if err != nil {
		slog.Error("Failed to get collection", "slug", collectionSlug, "error", err)
		if err == ErrNotFound {
			http.Error(w, "Collection not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	createdResource, err := h.Service.CreateResource(r.Context(), collection, slug, user.ID, fields)
	if err != nil {
		slog.Error("Failed to create resource", "collectionSlug", collectionSlug, "resourceSlug", slug, "error", err)
		if err == ErrAlreadyExists {
			http.Error(w, "Resource with this slug already exists", http.StatusConflict)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	util.JSON(w, http.StatusCreated, createdResource)
}

func (h *Handler) GetResource(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	resourceSlug := r.PathValue("resourceSlug")

	collection, err := h.Service.FindBySlug(r.Context(), slug)
	if err != nil {
		slog.Error("Failed to get collection", "slug", slug, "error", err)
		if err == ErrNotFound {
			http.Error(w, "Collection not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	if collection.IsGlobal {
		resourceSlug = slug
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

type FindAllQueryString struct {
	Search string `query:"q"`
}

func (h *Handler) FindAll(w http.ResponseWriter, r *http.Request) {
	query, err := util.QueryString[FindAllQueryString](r)
	if err != nil {
		slog.Error("Failed to decode query parameters", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	collections, err := h.Service.FindAll(r.Context(), &FindAllParams{Search: query.Search})
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

type FindAllGlobalsQueryString struct {
	Search string `query:"q"`
}

func (h *Handler) FindAllGlobals(w http.ResponseWriter, r *http.Request) {
	query, err := util.QueryString[FindAllQueryString](r)
	if err != nil {
		slog.Error("Failed to decode query parameters", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	globals, err := h.Service.FindAllGlobals(r.Context(), &FindAllParams{Search: query.Search})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := make([]CollectionResponse, len(globals))
	for i, global := range globals {
		response[i] = *NewCollectionResponse(&global)
	}

	util.JSON(w, http.StatusOK, response)
}

func (h *Handler) DeleteResource(w http.ResponseWriter, r *http.Request) {
	user := auth.RequestUser(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	slug := r.PathValue("slug")
	resourceSlug := r.PathValue("resourceSlug")

	collection, err := h.Service.FindBySlug(r.Context(), slug)
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

	if err := h.Service.DeleteResource(r.Context(), resource); err != nil {
		slog.Error("Failed to delete resource", "slug", slug, "resourceSlug", resourceSlug, "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
