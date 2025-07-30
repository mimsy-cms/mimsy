package collection

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	Service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) Definition(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("collectionSlug")
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

func (h *Handler) Items(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("collectionSlug")
	if slug == "" {
		http.Error(w, "Missing slug", http.StatusBadRequest)
		return
	}

	items, err := h.Service.GetItems(r.Context(), slug)
	if err != nil {
		if err == ErrNotFound {
			http.Error(w, "Collection not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	collections, err := h.Service.ListCollections(r.Context())
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(collections)
}
