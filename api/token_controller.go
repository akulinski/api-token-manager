package api

import (
	"encoding/json"
	"github.com/akulinski/api-token-manager/db"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"time"
)

var _ = godotenv.Load()

type TokenModel struct {
	Token       string    `json:"token"`
	GeneratedAt time.Time `json:"generatedAt"`
}

func AddToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var token db.Token

	err := json.NewDecoder(r.Body).Decode(&token)

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := db.Insert(token)

	err = json.NewEncoder(w).Encode(&result)

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}

func ValidateToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var tokenModel TokenModel

	err := json.NewDecoder(r.Body).Decode(&tokenModel)
	if err != nil {
		log.Println(err)
	}
	tokenStr, err := ValidateJwt(tokenModel)

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}

	if !tokenStr.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func GetAllTokens(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(db.GetAll())
}

func GetTokenById(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	id := getIdFromRequest(r)

	byId := db.GetById(id)

	json.NewEncoder(w).Encode(&byId)
}

func RevokeTokenApi(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := getIdFromRequest(r)

	revoked := db.RevokeToken(id)

	json.NewEncoder(w).Encode(&revoked)

}

func GenerateTokenForUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	username := getParamFromRequest(r, "username")
	tokenString := GenerateToken(username)

	tokenResponse := TokenModel{Token: tokenString, GeneratedAt: time.Now()}

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
