package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct{}

func NewAuthService() AuthService {
	return AuthService{}
}

type Claims struct {
	Sub int       `json:"sub"`
	Exp time.Time `json:"exp"`
	jwt.RegisteredClaims
}

func (s *AuthService) CreateToken(claims jwt.Claims, method jwt.SigningMethod, key []byte) (string, error) {
	token := jwt.NewWithClaims(method, claims)
	tokenString, err := token.SignedString(key)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) DecodeToken(tokenString string, method jwt.SigningMethod, key []byte) (*Claims, error) {
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
