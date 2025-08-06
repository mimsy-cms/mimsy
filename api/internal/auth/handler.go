package auth

import (
	"encoding/json"
	"net/http"

	"github.com/mimsy-cms/mimsy/internal/util"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginHandler(s AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := util.DecodeJSON[LoginRequest](r)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		resp, err := s.Login(r.Context(), req.Email, req.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		util.JSON(w, http.StatusOK, resp)
	}
}

func LogoutHandler(s AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			http.Error(w, "No session found", http.StatusUnauthorized)
			return
		}

		if err := s.Logout(r.Context(), cookie.Value); err != nil {
			http.Error(w, "Failed to delete session", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully", "session": cookie.Value})
	}
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func ChangePasswordHandler(s AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		if err := s.ChangePassword(r.Context(), user.ID, req.OldPassword, req.NewPassword); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		util.JSON(w, http.StatusOK, struct{}{})
	}
}

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"isAdmin"`
}

func RegisterHandler(s AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := util.DecodeJSON[CreateUserRequest](r)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if err := s.Register(r.Context(), *req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		util.JSON(w, http.StatusCreated, struct{}{})
	}
}

type MeResponse struct {
	ID                 int64  `json:"id"`
	Email              string `json:"email"`
	IsAdmin            bool   `json:"is_admin"`
	MustChangePassword bool   `json:"must_change_password"`
}

func MeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
}
