package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type UserHandler struct {
	db *sql.DB
}

type UserResponse struct {
	ID    int    `json:"id"`
	Nome  string `json:"nome"`
	Idade int    `json:"idade"`
}

type UserRequest struct {
	Nome  string `json:"nome"`
	Idade int    `json:"idade"`
}

func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{db: db}
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
	rows, err := h.db.Query("SELECT id, nome, idade FROM usuarios;")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	defer rows.Close()

	var data []UserResponse
	var userResponse UserResponse

	for rows.Next() {
		rows.Scan(&userResponse.ID, &userResponse.Nome, &userResponse.Idade)
		data = append(data, userResponse)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (h *UserHandler) GetOneUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.Path)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		return
	}

	var user UserResponse

	row := h.db.QueryRow("SELECT id, nome, idade FROM usuarios WHERE id = ?", id)
	if err := row.Scan(&user.ID, &user.Nome, &user.Idade); err != nil {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": fmt.Sprintf("User %d not found", id)})
		return
	}

	respondJSON(w, http.StatusOK, user)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser UserRequest
	json.NewDecoder(r.Body).Decode(&newUser)

	var tmp any
	row := h.db.QueryRow("SELECT nome FROM usuarios WHERE nome = ?", newUser.Nome)
	if err := row.Scan(&tmp); err == nil {
		respondJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": fmt.Sprintf("User with name %s already exists", newUser.Nome)})
		return
	}

	_, err := h.db.Exec("INSERT INTO usuarios(nome, idade) VALUES (?, ?)", newUser.Nome, newUser.Idade)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, newUser)
}
