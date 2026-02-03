package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/Cinnamoon-dev/blue-gopher/internal/database"
	"github.com/Cinnamoon-dev/blue-gopher/internal/http/handlers"
	"github.com/Cinnamoon-dev/blue-gopher/internal/middleware"
	"github.com/Cinnamoon-dev/blue-gopher/internal/repositories"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	mux := http.NewServeMux()
	db, err := sql.Open("sqlite3", "./storage.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	database.CreateTables("../internal/database/tables.sql", db)
	database.Populate("../internal/database/rules.sql", db)

	userRepository := repositories.NewUserRepository(db)
	userHandler := handlers.NewUserHandler(userRepository)
	authHandler := handlers.NewAuthHandler(userRepository)

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
	mux.Handle("/user/", middleware.Logging(middleware.Auth("user", userRepository, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	}))))

	mux.Handle("/auth", middleware.Logging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			authHandler.Login(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})))

	PORT, unset := os.LookupEnv("PORT")
	if unset == false {
		PORT = "3001"
	}

	log.Printf("Listening on port %s\n", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, mux))
}
