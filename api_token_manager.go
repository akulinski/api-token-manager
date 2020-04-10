package main

import (
	"fmt"
	"github.com/akulinski/api-token-manager/api"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

var _ = godotenv.Load()


func main() {

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	r := mux.NewRouter()

	r.HandleFunc("/api/v1/token", api.AddToken).Methods("POST")
	r.HandleFunc("/api/v1/token", api.GetAllTokens).Methods("GET")
	r.HandleFunc("/api/v1/token/{id}", api.GetTokenById).Methods("GET")
	r.HandleFunc("/api/v1/token/{id}/revoke", api.RevokeTokenApi).Methods("PATCH")
	r.HandleFunc("/api/v1/token/generate/{username}", api.GenerateTokenForUser).Methods("POST")
	r.HandleFunc("/api/v1/token/validate", api.ValidateToken).Methods("POST")

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}
