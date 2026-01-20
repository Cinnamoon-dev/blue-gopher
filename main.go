package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Cinnamoon-dev/blue-gopher/handlers"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	mux := http.NewServeMux()
	db, err := sql.Open("sqlite3", "./storage.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userHandler := handlers.NewUserHandler(db)

	// Expected URL: /user
	mux.Handle("/user", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet: // List all users
			userHandler.GetAllUsers(w, r)
		case http.MethodPost: // Create an user
			userHandler.CreateUser(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}))

	// Expected URL: /user/{id}
	mux.Handle("/user/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	}))

	log.Fatal(http.ListenAndServe(":3001", mux))
}
