package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/Cinnamoon-dev/blue-gopher/internal/http/handlers"
	"github.com/Cinnamoon-dev/blue-gopher/internal/repositories"
	"github.com/Cinnamoon-dev/blue-gopher/internal/services"
	"github.com/Cinnamoon-dev/blue-gopher/pkg/config"
	"github.com/golang-jwt/jwt/v5"
)

func Auth(controller string, repo repositories.UserRepository, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
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

		if claims.Exp.Before(time.Now()) {
			handlers.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "Token expired"})
			return
		}

		id := claims.Sub
		user, err := repo.Get(ctx, id)
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

		perms, err := repo.GetPermission(ctx, user.ID, action[r.Method], controller)
		if err != nil {
			handlers.RespondError(w, err)
			return
		}

		if perms == false {
			handlers.RespondJSON(w, http.StatusUnauthorized, map[string]string{"error": "User does not have permission"})
			return
		}

		ctx = context.WithValue(r.Context(), config.UserContextKey, user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
