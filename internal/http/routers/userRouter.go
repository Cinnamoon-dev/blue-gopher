package routers

import (
	"net/http"

	"github.com/Cinnamoon-dev/blue-gopher/internal/http/handlers"
	"github.com/Cinnamoon-dev/blue-gopher/internal/middleware"
	"github.com/Cinnamoon-dev/blue-gopher/pkg/config"
)

type UserRouter struct {
	UserHandler handlers.UserHandler
}

func NewUserRouter(userHandler handlers.UserHandler) UserRouter {
	return UserRouter{UserHandler: userHandler}
}

// Expected URL: /user
func (ro *UserRouter) BaseRoutes() http.HandlerFunc {
	router := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet: // List all users
			ro.UserHandler.GetAllUsers(w, r)
		case http.MethodPost: // Create an user
			ro.UserHandler.CreateUser(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	timeoutRouter := middleware.Timeout(config.DefaultTimeout, router)
	loggedTimeoutRouter := middleware.Logging(timeoutRouter)
	return loggedTimeoutRouter
}

// Expected URL: /user/{id}
func (ro *UserRouter) IDRoutes() http.HandlerFunc {
	router := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet: // Get one user with id
			ro.UserHandler.GetOneUser(w, r)
		case http.MethodPut: // Edit one user with id
			ro.UserHandler.UpdateUser(w, r)
		case http.MethodDelete: // Delete one user with id
			ro.UserHandler.DeleteUser(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	authRouter := middleware.Auth("user", ro.UserHandler.Svc.UserRepo, router)
	timeoutAuthRouter := middleware.Timeout(config.DefaultTimeout, authRouter)
	loggedTimeoutAuthRouter := middleware.Logging(timeoutAuthRouter)
	return loggedTimeoutAuthRouter
}
