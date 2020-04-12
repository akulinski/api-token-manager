package api

import (
	"encoding/json"
	"github.com/akulinski/api-token-manager/db"
	"github.com/akulinski/api-token-manager/domain"
	"github.com/akulinski/api-token-manager/services"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"time"
)

var _ = godotenv.Load()

var tokenRepository = db.NewTokenRepository()

func AddToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var token domain.Token

	err := json.NewDecoder(r.Body).Decode(&token)

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := tokenRepository.Insert(token)

	err = json.NewEncoder(w).Encode(&result)

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}

func ValidateToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tokenModel := getTokenModelFromRequest(r)

	token := tokenRepository.FindByTokenValue(tokenModel.Token)

	if token.Revoked == true {

		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	tokenStr, err := services.ValidateJwt(tokenModel)

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

	json.NewEncoder(w).Encode(tokenRepository.GetAll())
}

func GetTokenById(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	id := getIdFromRequest(r)

	byId := tokenRepository.GetById(id)

	json.NewEncoder(w).Encode(&byId)
}

func GetTokenByModel(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	tokenModel := getTokenModelFromRequest(r)

	token := tokenRepository.FindByTokenValue(tokenModel.Token)

	json.NewEncoder(w).Encode(&token)
}

func RevokeTokenApi(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tokenModel := getTokenModelFromRequest(r)

	token := tokenRepository.FindByTokenValue(tokenModel.Token)

	if token.Id == "" {
		log.Println("Failed to find token by value")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := primitive.ObjectIDFromHex(token.Id)

	if err != nil {
		log.Println("Failed to convert id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	revoked := tokenRepository.RevokeToken(id)

	json.NewEncoder(w).Encode(&revoked)
}

func GenerateTokenForUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var userModel domain.UserModel

	err := json.NewDecoder(r.Body).Decode(&userModel)
	if err != nil {
		log.Println(err)
	}

	tokenString := services.GenerateToken(userModel.UserIdentificator)

	tokenResponse := domain.TokenModel{Token: tokenString, GeneratedAt: time.Now()}

	token := domain.Token{
		IssuedAt: time.Now(),
		Issuer:   "SYSTEM",
		UserID:   userModel.UserIdentificator,
		Token:    tokenString,
		Expired:  false,
		Revoked:  false,
	}

	tokenRepository.Insert(token)

	json.NewEncoder(w).Encode(&tokenResponse)
}

func getIdFromRequest(r *http.Request) primitive.ObjectID {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])

	if err != nil {
		log.Println(err)
	}

	return id
}

func getTokenModelFromRequest(r *http.Request) domain.TokenModel {
	var tokenModel domain.TokenModel

	err := json.NewDecoder(r.Body).Decode(&tokenModel)
	if err != nil {
		log.Println(err)
	}

	return tokenModel
}
