package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Cinnamoon-dev/blue-gopher/internal/domain"
	"github.com/Cinnamoon-dev/blue-gopher/internal/services"
)

type UserHandler struct {
	Svc services.UserService
}

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	RoleID   int64  `json:"role_id"`
}

func NewUserHandler(svc services.UserService) UserHandler {
	return UserHandler{Svc: svc}
}

func respondJSON(w http.ResponseWriter, status int, payload any) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")

	if payload != nil {
		json.NewEncoder(w).Encode(payload)
	}
}

// The idea is simple: each handler is going to have its own parseID
// So each handler can parse the URL the way they want
func parseID(path string) (int, error) {
	// path = /user/{id}
	parts := strings.Split(path, "/")

	if len(parts) < 3 {
		return 0, http.ErrNotSupported
	}

	return strconv.Atoi(parts[2])
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.Svc.GetAll()
	if err != nil {
		switch err.Error() {
		case "Database error":
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		default:
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
		}
		return
	}

	respondJSON(w, http.StatusOK, users)
}

func (h *UserHandler) GetOneUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.Path)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		return
	}

	user, err := h.Svc.Get(id)
	if err != nil {
		switch err.Error() {
		case "Not found":
			respondJSON(w, http.StatusNotFound, map[string]string{"error": fmt.Sprintf("User %d not found", id)})
		case "Database error":
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		default:
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
		}
		return
	}

	respondJSON(w, http.StatusOK, user)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser UserRequest
	json.NewDecoder(r.Body).Decode(&newUser)
	newUser.Email = strings.TrimSpace(newUser.Email)

	user := domain.User{
		ID:       0,
		Email:    newUser.Email,
		Password: newUser.Password,
		RoleID:   newUser.RoleID,
	}

	if err := user.ValidateEmail(); err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	if err := user.ValidatePassword(); err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	id, err := h.Svc.Create(user)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("User %d created successfully", id)})
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.Path)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		return
	}

	var fields UserRequest
	if err := json.NewDecoder(r.Body).Decode(&fields); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		return
	}

	user := domain.User{
		ID:       0,
		Email:    fields.Email,
		Password: fields.Password,
		RoleID:   fields.RoleID,
	}

	if err := user.ValidateEmail(); err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	if err := user.ValidatePassword(); err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	if err := h.Svc.Update(id, user); err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondJSON(w, http.StatusOK, fields)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.Path)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		return
	}

	if err := h.Svc.Delete(id); err != nil {
		switch err.Error() {
		case "Database error":
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		default:
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
		}
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("User %d deleted successfully", id)})
}
