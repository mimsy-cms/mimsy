package sync

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/mimsy-cms/mimsy/internal/auth"
	"github.com/mimsy-cms/mimsy/internal/cron"
	"github.com/mimsy-cms/mimsy/internal/util"
)

type Handler struct {
	Repository  SyncStatusRepository
	CronService cron.CronService
}

func NewHandler(repository SyncStatusRepository, cronService cron.CronService) *Handler {
	return &Handler{
		Repository:  repository,
		CronService: cronService,
	}
}

type StatusQueryString struct {
	Limit int `query:"limit"`
}

func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {
	user := auth.RequestUser(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	query, err := util.QueryString[StatusQueryString](r)
	if err != nil {
		slog.Error("Failed to decode query parameters", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	limit := query.Limit
	if limit <= 0 || limit > 10 {
		limit = 5 // Default to 5, max 10
	}

	statuses, err := h.Repository.GetRecentStatuses(limit)
	if err != nil {
		slog.Error("Failed to get recent sync statuses", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	util.JSON(w, http.StatusOK, NewStatusResponse(statuses, os.Getenv("GH_REPO")))
}

func (h *Handler) Jobs(w http.ResponseWriter, r *http.Request) {
	user := auth.RequestUser(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	jobStatuses, err := h.CronService.GetJobStatuses(r.Context())
	if err != nil {
		slog.Error("Failed to get job statuses", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	util.JSON(w, http.StatusOK, jobStatuses)
}

func (h *Handler) ActiveMigration(w http.ResponseWriter, r *http.Request) {
	user := auth.RequestUser(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	repo := os.Getenv("GH_REPO")
	if repo == "" {
		slog.Error("GH_REPO environment variable not set")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	activeMigration, err := h.Repository.GetActiveMigration(repo)
	if err != nil {
		slog.Error("Failed to get active migration", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if activeMigration == nil {
		util.JSON(w, http.StatusOK, map[string]interface{}{"active_migration": nil})
		return
	}

	util.JSON(w, http.StatusOK, map[string]interface{}{"active_migration": activeMigration})
}
