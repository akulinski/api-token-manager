package api

import (
	"github.com/dgrijalva/jwt-go"
	"log"
	"os"
	"time"
)

var jwtKey = os.Getenv("JWT_SECRET")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateToken(username string) string {

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString([]byte(jwtKey))

	if err != nil {
		log.Printf("Failed to generate jwt token %s", err)
	}

	return tokenString
}

func ValidateJwt(model TokenModel) (*jwt.Token, error) {

	claims := Claims{}

	tkn, err := jwt.ParseWithClaims(model.Token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})

	return tkn, err
}