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

	setUpV1Routes(r)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}

func setUpV1Routes(r *mux.Router) {
	r.HandleFunc("/api/v1/token", api.AddToken).Methods("POST")
	r.HandleFunc("/api/v1/token", api.GetAllTokens).Methods("GET")
	r.HandleFunc("/api/v1/token/find", api.GetTokenByModel).Methods("POST")
	r.HandleFunc("/api/v1/token/find/{id}", api.GetTokenById).Methods("GET")

	r.HandleFunc("/api/v1/token/revoke", api.RevokeTokenApi).Methods("PATCH")
	r.HandleFunc("/api/v1/token/generate", api.GenerateTokenForUser).Methods("POST")
	r.HandleFunc("/api/v1/token/validate", api.ValidateToken).Methods("POST")
}
