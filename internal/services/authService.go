package services

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Cinnamoon-dev/blue-gopher/internal/customerrors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct{}

func NewAuthService() AuthService {
	return AuthService{}
}

type AuthClaims struct {
	Sub int64     `json:"sub"`
	Exp time.Time `json:"exp"`
	jwt.RegisteredClaims
}

type MailClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func DecodeToken[T jwt.Claims](tokenString string, method jwt.SigningMethod, key []byte, claims T) (T, error) {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != method.Alg() {
			return claims, &customerrors.HTTPError{Status: http.StatusBadRequest, Message: fmt.Sprintf("Invalid algorithm: %s", token.Method.Alg())}
		}

		return key, nil
	})

	if err != nil {
		return claims, &customerrors.HTTPError{Status: http.StatusBadRequest, Message: err.Error()}
	}

	if claims, ok := token.Claims.(T); ok {
		return claims, nil
	}

	return claims, &customerrors.HTTPError{Status: http.StatusBadRequest, Message: "unknown claims type, cannot proceed"}
}

func (s *AuthService) CreateToken(claims jwt.Claims, method jwt.SigningMethod, key []byte) (string, error) {
	token := jwt.NewWithClaims(method, claims)
	tokenString, err := token.SignedString(key)

	if err != nil {
		return "", &customerrors.HTTPError{Status: http.StatusInternalServerError, Message: err.Error()}
	}

	return tokenString, nil
}

func (s *AuthService) DecodeAuthToken(tokenString string, method jwt.SigningMethod, key []byte) (*AuthClaims, error) {
	return DecodeToken(tokenString, method, key, &AuthClaims{})
}

func (s *AuthService) DecodeMailToken(tokenString string, method jwt.SigningMethod, key []byte) (*MailClaims, error) {
	return DecodeToken(tokenString, method, key, &MailClaims{})
}

func (s *AuthService) Hash(text string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (s *AuthService) VerifyHash(text string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(text))
	return err == nil
}
