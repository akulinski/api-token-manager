package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

var err = godotenv.Load()

func main() {

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	r := mux.NewRouter()

	r.HandleFunc("/api/v1/token", AddToken).Methods("POST")
	r.HandleFunc("/api/v1/token", GetAllTokens).Methods("GET")
	r.HandleFunc("/api/v1/token/{id}", GetTokenById).Methods("GET")
	r.HandleFunc("/api/v1/token/{id}/revoke", RevokeTokenApi).Methods("PATCH")

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))

}
