package database

import (
	"database/sql"

	"github.com/Cinnamoon-dev/blue-gopher/pkg/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "modernc.org/sqlite"
)

func RunAllMigrations(db *sql.DB) {
	env := config.NewEnv()
	m, err := migrate.New(env.MigrationsUrl, "sqlite3://"+env.DbUrl)
	if err != nil {
		panic(err)
	}
	defer m.Close()

	err = m.Up()
	if err != nil && err.Error() != "no change" {
		panic(err)
	}
}

func UndoAllMigrations(db *sql.DB) {
	env := config.NewEnv()
	m, err := migrate.New(env.MigrationsUrl, "sqlite3://"+env.DbUrl)
	if err != nil {
		panic(err)
	}
	defer m.Close()

	err = m.Down()
	if err != nil && err.Error() != "no change" {
		panic(err)
	}
}
