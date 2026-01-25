package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Cinnamoon-dev/blue-gopher/repositories"
)

type UserHandler struct {
	Repo repositories.UserRepository
}

type UserRequest struct {
	Nome  string `json:"nome"`
	Idade int    `json:"idade"`
}

func NewUserHandler(Repo repositories.UserRepository) *UserHandler {
	return &UserHandler{Repo: Repo}
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
	users, err := h.Repo.GetAll()
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

	user, err := h.Repo.Get(id)
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
	var newUser repositories.User
	json.NewDecoder(r.Body).Decode(&newUser)

	_, err := h.Repo.GetByName(newUser.Nome)
	if err == nil {
		respondJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": fmt.Sprintf("User with name %s already exists", newUser.Nome)})
		return
	}

	newUser.Nome = strings.TrimSpace(newUser.Nome)

	id, err := h.Repo.Create(newUser)
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

	var fields repositories.User
	if err := json.NewDecoder(r.Body).Decode(&fields); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		return
	}
	fmt.Printf("%+v\n\n", fields)

	if _, err := h.Repo.Get(id); err != nil {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": fmt.Sprintf("User %d not found", id)})
		return
	}

	if fields.Nome == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Field 'nome' is required"})
		return
	}

	if fields.Idade < 1 {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Field 'idade' should be greater than 0"})
		return
	}

	if err := h.Repo.Update(id, fields); err != nil {
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

	if err := h.Repo.Delete(id); err != nil {
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
