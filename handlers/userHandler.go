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

func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{db: db}
}

func respondJSON(w http.ResponseWriter, status int, payload any) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	if payload != nil {
		json.NewEncoder(w).Encode(payload)
	}
}

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
	json.NewEncoder(w).Encode(map[string]any{"data": data})
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
