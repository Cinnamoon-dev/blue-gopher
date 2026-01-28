package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Cinnamoon-dev/blue-gopher/handlers"
	"github.com/Cinnamoon-dev/blue-gopher/middleware"
	"github.com/Cinnamoon-dev/blue-gopher/repositories"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	mux := http.NewServeMux()
	db, err := sql.Open("sqlite3", "./storage.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	PORT, unset := os.LookupEnv("PORT")
	if unset == false {
		PORT = "3001"
	}

	userRepository := repositories.NewUserRepository(db)
	userHandler := handlers.NewUserHandler(userRepository)

	// Expected URL: /user
	mux.Handle("/user", middleware.Logging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet: // List all users
			userHandler.GetAllUsers(w, r)
		case http.MethodPost: // Create an user
			userHandler.CreateUser(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})))

	// Expected URL: /user/{id}
	mux.Handle("/user/", middleware.Logging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method { // URL: /user/{id}
		case http.MethodGet: // Get one user with id
			userHandler.GetOneUser(w, r)
		case http.MethodPut: // Edit one user with id
			userHandler.UpdateUser(w, r)
		case http.MethodDelete: // Delete one user with id
			userHandler.DeleteUser(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})))

	fmt.Printf("Listening on port %s\n", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, mux))
}
