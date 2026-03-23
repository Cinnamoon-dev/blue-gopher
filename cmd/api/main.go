package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Cinnamoon-dev/blue-gopher/internal/database"
	"github.com/Cinnamoon-dev/blue-gopher/internal/http/handlers"
	"github.com/Cinnamoon-dev/blue-gopher/internal/http/routers"
	"github.com/Cinnamoon-dev/blue-gopher/internal/messaging/rabbitmq"
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
	db.Exec("PRAGMA jounal_mode=WAL;")

	database.CreateTables("./internal/database/tables.sql", db)
	database.Populate("./internal/database/rules.sql", db)
	database.RunAllMigrations(db)

	queueConn, err := rabbitmq.NewConnection(env.RabbitMQUrl)
	if err != nil {
		log.Panicf("Couldn't connect to RabbitMQ: %s", err)
	}

	publisher := rabbitmq.NewRabbitPublisher(*queueConn)

	roleRepository := repositories.NewRoleRepository(db)

	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository, roleRepository, *publisher)
	userHandler := handlers.NewUserHandler(userService)
	userRouter := routers.NewUserRouter(userHandler)

	authService := services.NewAuthService()
	authHandler := handlers.NewAuthHandler(userRepository)
	authRouter := routers.NewAuthRouter(authHandler)

	mailService := services.NewMailService()
	mailHandler := handlers.NewMailHandler(authService, userService, mailService)
	mailRouter := routers.NewMailRouter(mailHandler)

	mux.Handle("/user", userRouter.BaseRoutes())
	mux.Handle("/user/", userRouter.IDRoutes())

	mux.Handle("/auth", authRouter.BaseRoutes())

	mux.Handle("/mail/", mailRouter.BaseRoutes())

	log.Printf("Listening on port %s\n", env.Port)
	log.Fatal(http.ListenAndServe(":"+env.Port, mux))
}
