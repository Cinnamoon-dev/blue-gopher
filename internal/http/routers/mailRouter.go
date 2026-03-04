package routers

import (
	"net/http"
	"strings"

	"github.com/Cinnamoon-dev/blue-gopher/internal/http/handlers"
	"github.com/Cinnamoon-dev/blue-gopher/internal/middleware"
	"github.com/Cinnamoon-dev/blue-gopher/internal/services"
	"github.com/Cinnamoon-dev/blue-gopher/pkg/config"
	"github.com/golang-jwt/jwt/v5"
)

type MailRouter struct {
	AuthService services.AuthService
	UserService services.UserService
}

func NewMailRouter(authService services.AuthService, userService services.UserService) MailRouter {
	return MailRouter{
		AuthService: authService,
		UserService: userService,
	}
}

func parseToken(path string) (string, error) {
	// path = /mail/{token}
	parts := strings.Split(path, "/")

	if len(parts) < 3 {
		return "", http.ErrNotSupported
	}

	return parts[2], nil
}

func (ro *MailRouter) BaseRoutes() http.HandlerFunc {
	router := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		switch r.Method {
		case http.MethodGet:
			env := config.NewEnv()
			emailToken, err := parseToken(r.URL.Path)
			if err != nil {
				handlers.RespondError(w, err)
				return
			}

			claims, err := ro.AuthService.DecodeMailToken(emailToken, jwt.SigningMethodHS256, []byte(env.JwtKey))
			if err != nil {
				handlers.RespondError(w, err)
				return
			}

			email := claims.Email
			user, err := ro.UserService.GetByEmail(ctx, email)
			if err != nil {
				handlers.RespondError(w, err)
				return
			}

			handlers.RespondJSON(w, http.StatusOK, user)
			return
			// Adicionar a coluna isverified a user table e mexer com todos os metodos de user repository
			// adicionar o datetime atual ao is_verified se for nulo se ja for verificado, não fazer nada e dizer que ja esta verificado
		}
	})

	LoggingRouter := middleware.Logging(router)
	return LoggingRouter
}
