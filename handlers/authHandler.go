package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Cinnamoon-dev/blue-gopher/middleware"
	"github.com/Cinnamoon-dev/blue-gopher/repositories"
	"github.com/golang-jwt/jwt/v5"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthHandler struct {
	Repo repositories.UserRepository
}

func NewAuthHandler(Repo repositories.UserRepository) AuthHandler {
	return AuthHandler{Repo: Repo}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var request LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	request.Email = strings.TrimSpace(request.Email)
	request.Password = strings.TrimSpace(request.Password)

	if request.Email == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "field email is required"})
		return
	}

	if request.Password == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "field password is required"})
		return
	}

	user, err := h.Repo.GetByEmail(request.Email)
	if err != nil {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "Email not found"})
		return
	}

	// TODO
	// Password hash
	if user.Password != request.Password {
		respondJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": "Wrong Password"})
		return
	}

	key := os.Getenv("SECRET_KEY")
	if key == "" {
		key = "d0699dddcf3e6896ff556dc156a6d65931a855b327822dc12ea5f67350125a45"
	}

	accessToken, err := middleware.CreateToken(
		jwt.MapClaims{
			"sub": user.ID,
			"exp": time.Now().Add(20 * time.Minute),
		},
		jwt.SigningMethodHS256,
		[]byte(key),
	)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	refreshToken, err := middleware.CreateToken(
		jwt.MapClaims{
			"sub": user.ID,
			"exp": time.Now().Add(7 * 24 * time.Hour),
		},
		jwt.SigningMethodHS256,
		[]byte(key),
	)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}
