package collection

import (
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
			http.Error(w, "Collection not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	util.JSON(w, http.StatusOK, resources)
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

	util.JSON(w, http.StatusOK, globals)
}
