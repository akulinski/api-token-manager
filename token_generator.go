package main

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
	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	// Create the JWT string
	tokenString, err := token.SignedString([]byte(jwtKey))

	if err!=nil{
		log.Printf("Failed to generate jwt token %s", err)
	}

	return tokenString
}
