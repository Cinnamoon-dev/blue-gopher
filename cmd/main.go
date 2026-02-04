package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Cinnamoon-dev/blue-gopher/internal/database"
	"github.com/Cinnamoon-dev/blue-gopher/internal/http/handlers"
	"github.com/Cinnamoon-dev/blue-gopher/internal/http/routers"
	"github.com/Cinnamoon-dev/blue-gopher/internal/repositories"
	"github.com/Cinnamoon-dev/blue-gopher/internal/services"
	"github.com/Cinnamoon-dev/blue-gopher/pkg/config"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	mux := http.NewServeMux()
	db, err := sql.Open("sqlite3", "./storage.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	env := config.NewEnv()

	database.CreateTables("../internal/database/tables.sql", db)
	database.Populate("../internal/database/rules.sql", db)

	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)
	userRouter := routers.NewUserRouter(userHandler)

	authHandler := handlers.NewAuthHandler(userRepository)
	authRouter := routers.NewAuthRouter(authHandler)

	mux.Handle("/user", userRouter.BaseRoutes())
	mux.Handle("/user/", userRouter.IDRoutes())

	mux.Handle("/auth", authRouter.BaseRoutes())

	log.Printf("Listening on port %s\n", env.Port)
	log.Fatal(http.ListenAndServe(":"+env.Port, mux))
}
