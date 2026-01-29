package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type MyClaims struct {
	Sub int       `json:"sub"`
	Exp time.Time `json:"exp"`
	jwt.RegisteredClaims
}

func create_token(claims jwt.Claims, method jwt.SigningMethod, key []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(key)

	if err != nil {
		return "", err
	}

	return tokenString, err
}

func decode_token(tokenString string, method jwt.SigningMethod, key string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("Invalid algorithm: %s", token.Method.Alg())
		}

		return []byte(key), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*MyClaims); ok {
		return claims, nil
	}

	return nil, errors.New("unknown claims type, cannot proceed")
}

func main() {
	key := []byte("480feb81fe9ba214177c977b521fcacc6a760cef6053782d43a699c534292766")
	type MyClaims struct {
		Sub int       `json:"sub"`
		Exp time.Time `json:"exp"`
		jwt.RegisteredClaims
	}

	// JWT with claims
	var realToken *jwt.Token
	realToken = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": 1,
		"exp": time.Now(),
	})
	realTokenString, err1 := realToken.SignedString(key)
	if err1 != nil {
		log.Fatal(err1.Error())
	}
	fmt.Printf("real token: %s\n\n", realTokenString)

	// JWT Decode
	token, err := jwt.ParseWithClaims(realTokenString, &MyClaims{}, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("Invalid algorithm: %s", token.Method.Alg())
		}

		return []byte(key), nil
	})

	if err != nil {
		log.Fatal(err.Error())
	} else if claims, ok := token.Claims.(*MyClaims); ok {
		fmt.Printf("sub: %d, exp: %s", claims.Sub, claims.Exp.String())
	} else {
		log.Fatal("unknown claims type, cannot proceed")
	}
}
