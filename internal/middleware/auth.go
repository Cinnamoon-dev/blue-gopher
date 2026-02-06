package middleware

import (
	"net/http"

	"github.com/Cinnamoon-dev/blue-gopher/internal/http/handlers"
	"github.com/Cinnamoon-dev/blue-gopher/internal/repositories"
	"github.com/Cinnamoon-dev/blue-gopher/internal/services"
	"github.com/Cinnamoon-dev/blue-gopher/pkg/config"
	"github.com/golang-jwt/jwt/v5"
)

func Auth(controller string, repo repositories.UserRepository, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Bearer")
		if token == "" {
			handlers.RespondJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not Authenticated"})
			return
		}

		authService := services.NewAuthService()
		env := config.NewEnv()
		claims, err := authService.DecodeToken(token, jwt.SigningMethodHS256, []byte(env.JwtKey))
		if err != nil {
			handlers.RespondError(w, err)
			return
		}

		id := claims.Sub
		user, err := repo.Get(id)
		if err != nil {
			handlers.RespondError(w, err)
			return
		}

		action := map[string]string{
			"GET":    "all",
			"POST":   "add",
			"PUT":    "edit",
			"DELETE": "delete",
		}

		perms, err := repo.GetPermission(user.ID, action[r.Method], controller)
		if err != nil {
			handlers.RespondError(w, err)
			return
		}

		if perms == false {
			handlers.RespondJSON(w, http.StatusUnauthorized, map[string]string{"error": "User does not have permission"})
			return
		}

		next.ServeHTTP(w, r)
	})
}
