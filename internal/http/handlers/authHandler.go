package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Cinnamoon-dev/blue-gopher/internal/repositories"
	"github.com/Cinnamoon-dev/blue-gopher/internal/services"
	"github.com/Cinnamoon-dev/blue-gopher/pkg/config"
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
		RespondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	request.Email = strings.TrimSpace(request.Email)
	request.Password = strings.TrimSpace(request.Password)

	if request.Email == "" {
		RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "field email is required"})
		return
	}

	if request.Password == "" {
		RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "field password is required"})
		return
	}

	user, err := h.Repo.GetByEmail(request.Email)
	if err != nil {
		RespondJSON(w, http.StatusNotFound, map[string]string{"error": "Email not found"})
		return
	}

	// TODO:
	// Password hash
	if user.Password != request.Password {
		RespondJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": "Wrong Password"})
		return
	}

	env := config.NewEnv()
	authService := services.NewAuthService()
	accessToken, err := authService.CreateToken(
		jwt.MapClaims{
			"sub": user.ID,
			"exp": time.Now().Add(20 * time.Minute),
		},
		jwt.SigningMethodHS256,
		[]byte(env.JwtKey),
	)
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	refreshToken, err := authService.CreateToken(
		jwt.MapClaims{
			"sub": user.ID,
			"exp": time.Now().Add(7 * 24 * time.Hour),
		},
		jwt.SigningMethodHS256,
		[]byte(env.JwtKey),
	)
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	authService := services.NewAuthService()
	token := r.Header.Get("Bearer")
	env := config.NewEnv()

	claims, err := authService.DecodeToken(token, jwt.SigningMethodHS256, []byte(env.JwtKey))
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
		return
	}

	id := claims.Sub
	user, err := h.Repo.Get(id)
	if err != nil {
		switch err.Error() {
		case "Not found":
			RespondJSON(w, http.StatusNotFound, map[string]string{"error": fmt.Sprintf("User %d not found", id)})
		case "Database error":
			RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		default:
			RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
		}
	}

	RespondJSON(w, http.StatusOK, user)
}
