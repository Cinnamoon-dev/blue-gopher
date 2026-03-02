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
	router := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			ro.AuthHandler.Login(w, r)
		case http.MethodGet:
			ro.AuthHandler.Me(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	loggedRouter := middleware.Logging(router)
	return loggedRouter
}
