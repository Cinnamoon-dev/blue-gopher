package config

import "os"

type Env struct {
	Port   string
	JwtKey string
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

	return Env{
		Port:   port,
		JwtKey: key,
	}
}
