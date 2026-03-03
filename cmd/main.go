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
	_ "modernc.org/sqlite"
)

func main() {
	env := config.NewEnv()
	mux := http.NewServeMux()
	db, err := sql.Open("sqlite", env.DbUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.Exec("PRAGMA journal_mode=WAL")

	database.CreateTables("./internal/database/tables.sql", db)
	database.Populate("./internal/database/rules.sql", db)
	database.RunAllMigrations(db)

	roleRepository := repositories.NewRoleRepository(db)

	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository, roleRepository)
	userHandler := handlers.NewUserHandler(userService)
	userRouter := routers.NewUserRouter(userHandler)

	authHandler := handlers.NewAuthHandler(userRepository)
	authRouter := routers.NewAuthRouter(authHandler)

	orderRouter := routers.NewOrderRouter(db)

	mux.Handle("/user", userRouter.BaseRoutes())
	mux.Handle("/user/", userRouter.IDRoutes())

	mux.Handle("/auth", authRouter.BaseRoutes())

	mux.Handle("/order", orderRouter.BaseRoutes())

	log.Printf("Listening on port %s\n", env.Port)
	log.Fatal(http.ListenAndServe(":"+env.Port, mux))
}
