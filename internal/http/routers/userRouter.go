package routers

import (
	"net/http"

	"github.com/Cinnamoon-dev/blue-gopher/internal/http/handlers"
	"github.com/Cinnamoon-dev/blue-gopher/internal/middleware"
)

type UserRouter struct {
	UserHandler handlers.UserHandler
}

func NewUserRouter(userHandler handlers.UserHandler) UserRouter {
	return UserRouter{UserHandler: userHandler}
}

// Expected URL: /user
func (ro *UserRouter) BaseRoutes() http.HandlerFunc {
	return middleware.Logging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet: // List all users
			ro.UserHandler.GetAllUsers(w, r)
		case http.MethodPost: // Create an user
			ro.UserHandler.CreateUser(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}))
}

// Expected URL: /user/{id}
func (ro *UserRouter) IDRoutes() http.HandlerFunc {
	return middleware.Logging(middleware.Auth("user", ro.UserHandler.Svc.Repo, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	})))
}
