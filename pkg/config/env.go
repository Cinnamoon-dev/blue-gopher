package config

import (
	"os"
	"path/filepath"
)

type Env struct {
	Port          string
	JwtKey        string
	DbUrl         string
	MigrationsUrl string
	BackendUrl    string
}

func NewEnv() Env {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	key := os.Getenv("JWT_KEY")
	if key == "" {
		key = "d0699dddcf3e6896ff556dc156a6d65931a855b327822dc12ea5f67350125a45"
	}

	backendUrl := os.Getenv("BACKEND_URL")
	if backendUrl == "" {
		backendUrl = "http://localhost:3001"
	}

	// Takes the path to the binary executable
	// So the binary needs to run in the root of the project
	// Or at least in this specific folder organization
	// In case the environment variables are not declared
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	dir := filepath.Dir(exePath)

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		dbUrl = dir + "/storage.db"
	}

	migrationsUrl := "file://" + dir + "/internal/database/migrations"

	return Env{
		Port:          port,
		JwtKey:        key,
		DbUrl:         dbUrl,
		MigrationsUrl: migrationsUrl,
		BackendUrl:    backendUrl,
	}
}
