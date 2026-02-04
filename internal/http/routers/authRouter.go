package routers

import (
	"net/http"

	"github.com/Cinnamoon-dev/blue-gopher/internal/http/handlers"
	"github.com/Cinnamoon-dev/blue-gopher/internal/middleware"
)

type AuthRouter struct {
	AuthHandler handlers.AuthHandler
}

func NewAuthRouter(authHandler handlers.AuthHandler) AuthRouter {
	return AuthRouter{AuthHandler: authHandler}
}

// Expected URL: /auth
func (ro *AuthRouter) BaseRoutes() http.HandlerFunc {
	return middleware.Logging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			ro.AuthHandler.Login(w, r)
		case http.MethodGet:
			// TODO: get current user
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}))
}
