package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"time"
)

type TokenResponse struct {
	Token       string    `json:"value"`
	GeneratedAt time.Time `json:"generatedAt"`
}

func AddToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var token Token

	err = json.NewDecoder(r.Body).Decode(&token)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := Insert(token)

	json.NewEncoder(w).Encode(&result)
}

func GetAllTokens(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(GetAll())
}

func GetTokenById(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	id := getIdFromRequest(r)

	byId := GetById(id)

	json.NewEncoder(w).Encode(&byId)
}

func RevokeTokenApi(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := getIdFromRequest(r)

	revoked := RevokeToken(id)

	json.NewEncoder(w).Encode(&revoked)

}

func GenerateTokenForUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	username := getParamFromRequest(r, "username")
	tokenString := GenerateToken(username)

	tokenResponse := TokenResponse{Token:tokenString, GeneratedAt:time.Now()}

	json.NewEncoder(w).Encode(&tokenResponse)

}

func getIdFromRequest(r *http.Request) primitive.ObjectID {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])

	if err != nil {
		log.Println(err)
	}

	return id
}

func getParamFromRequest(r *http.Request, param string) string {

	id := mux.Vars(r)[param]

	return id
}
