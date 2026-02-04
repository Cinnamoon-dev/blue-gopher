package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/Cinnamoon-dev/blue-gopher/internal/repositories"
	"github.com/Cinnamoon-dev/blue-gopher/internal/services"
	"github.com/Cinnamoon-dev/blue-gopher/pkg/config"
	"github.com/golang-jwt/jwt/v5"
)

func Auth(controller string, repo repositories.UserRepository, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Bearer")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"error": "Not Authenticated"})
			return
		}

		authService := services.NewAuthService()
		env := config.NewEnv()
		claims, err := authService.DecodeToken(token, jwt.SigningMethodHS256, []byte(env.JwtKey))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		id := claims.Sub
		user, err := repo.Get(id)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
			return
		}

		action := map[string]string{
			"all":    "GET",
			"add":    "POST",
			"edit":   "PUT",
			"delete": "DELETE",
		}

		perms, err := repo.GetPermission(user.ID, action[r.Method], controller)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		if perms == false {
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"error": "User does not have permission"})
			return
		}

		next.ServeHTTP(w, r)
	})
}
