package auth

import (
	"encoding/json"
	"net/http"

	"github.com/mimsy-cms/mimsy/internal/util"
)

type handler struct {
	authService AuthService
}

func NewHandler(authService AuthService) *handler {
	return &handler{authService: authService}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *handler) Login(w http.ResponseWriter, r *http.Request) {
	req, err := util.DecodeJSON[LoginRequest](r)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := s.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		switch err.Error() {
		case "user not found":
			http.Error(w, "User not found", http.StatusUnauthorized)
		case "invalid credentials":
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	util.JSON(w, http.StatusOK, resp)
}

func (s *handler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Error(w, "No session found", http.StatusUnauthorized)
		return
	}

	if err := s.authService.Logout(r.Context(), cookie.Value); err != nil {
		http.Error(w, "Failed to delete session", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully", "session": cookie.Value})
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (s *handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	req, err := util.DecodeJSON[ChangePasswordRequest](r)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = s.authService.ChangePassword(r.Context(), user.ID, req.OldPassword, req.NewPassword)
	if err != nil {
		switch err.Error() {
		case "old password is incorrect":
			http.Error(w, err.Error(), http.StatusUnauthorized)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	util.JSON(w, http.StatusOK, struct{}{})
}

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"isAdmin"`
}

func (s *handler) Register(w http.ResponseWriter, r *http.Request) {
	req, err := util.DecodeJSON[CreateUserRequest](r)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := s.authService.Register(r.Context(), *req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	util.JSON(w, http.StatusCreated, struct{}{})
}

type MeResponse struct {
	ID                 int64  `json:"id"`
	Email              string `json:"email"`
	IsAdmin            bool   `json:"is_admin"`
	MustChangePassword bool   `json:"must_change_password"`
}

func (s *handler) Me(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	util.JSON(w, http.StatusOK, MeResponse{
		ID:                 user.ID,
		Email:              user.Email,
		IsAdmin:            user.IsAdmin,
		MustChangePassword: user.MustChangePassword,
	})
}
