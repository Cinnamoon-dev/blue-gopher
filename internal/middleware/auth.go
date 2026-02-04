package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Cinnamoon-dev/blue-gopher/internal/repositories"
	"github.com/Cinnamoon-dev/blue-gopher/pkg/config"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Sub int       `json:"sub"`
	Exp time.Time `json:"exp"`
	jwt.RegisteredClaims
}

func CreateToken(claims jwt.Claims, method jwt.SigningMethod, key []byte) (string, error) {
	token := jwt.NewWithClaims(method, claims)
	tokenString, err := token.SignedString(key)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func DecodeToken(tokenString string, method jwt.SigningMethod, key []byte) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != method.Alg() {
			return nil, fmt.Errorf("Invalid algorithm: %s", token.Method.Alg())
		}

		return key, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok {
		return claims, nil
	}

	return nil, errors.New("unknown claims type, cannot proceed")
}

func Auth(controller string, repo repositories.UserRepository, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Bearer")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"error": "Not Authenticated"})
			return
		}

		env := config.NewEnv()
		claims, err := DecodeToken(token, jwt.SigningMethodHS256, []byte(env.JwtKey))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		id := claims.Sub
		user, err := repo.Get(id)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
			return
		}

		action := map[string]string{
			"all":    "GET",
			"add":    "POST",
			"edit":   "PUT",
			"delete": "DELETE",
		}

		perms, err := repo.GetPermission(user.ID, action[r.Method], controller)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		if perms == false {
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"error": "User does not have permission"})
			return
		}

		next.ServeHTTP(w, r)
	})
}
