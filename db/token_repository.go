package db

import (
	"context"
	"github.com/akulinski/api-token-manager/domain"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

var _ = godotenv.Load()

func init() {
	mod := mongo.IndexModel{
		Keys: bson.M{
			"token": 1,
		}, Options: options.Index().SetUnique(true),
	}

	_, err := getCollection().Indexes().CreateOne(context.Background(), mod)

	if err != nil {
		log.Fatalf("failed to create index %s", err)
	}
}

type TokenRepository interface {
	getCollection() *mongo.Collection
	FindById(id primitive.ObjectID) domain.Token
	GetAll() []*domain.Token
	Insert(token domain.Token) *mongo.InsertOneResult
	GetById(id primitive.ObjectID) domain.Token
	RevokeToken(id primitive.ObjectID) int64
	FindByTokenValue(tokenValue string) domain.Token
}

type tokenRepository struct {
	client *mongo.Client
}

func NewTokenRepository() tokenRepository {
	client := getClient()
	return tokenRepository{client: client}
}

func (r tokenRepository) RevokeToken(id primitive.ObjectID) int64 {
	filter := bson.M{"_id": id}

	byId := r.GetById(id)

	byId.Revoked = true

	result, err := r.getCollection().UpdateOne(context.Background(), filter, bson.D{
		{
			"$set", bson.D{{"revoked", true}},
		},
	})

	if err != nil {
		log.Println(err)
	}

	return result.ModifiedCount
}

func (r tokenRepository) FindByTokenValue(tokenValue string) domain.Token {
	filter := bson.M{"token": tokenValue}

	var result domain.Token

	err := r.getCollection().FindOne(context.Background(), filter).Decode(&result)

	if err != nil {
		log.Println(err)
	}

	return result
}

func (r tokenRepository) GetById(id primitive.ObjectID) domain.Token {
	filter := bson.M{"_id": id}

	var result domain.Token

	err := r.getCollection().FindOne(context.Background(), filter).Decode(&result)

	if err != nil {
		log.Println(err)
	}

	return result
}

func (r tokenRepository) Insert(token domain.Token) *mongo.InsertOneResult {
	collection := r.getCollection()
	res, err := collection.InsertOne(context.Background(), token)

	if err != nil {
		log.Printf("Failed to insert token %s", err)
	}

	return res
}

func (r tokenRepository) FindById(id primitive.ObjectID) domain.Token {
	filter := bson.M{"_id": id}

	var result domain.Token

	err := r.getCollection().FindOne(context.Background(), filter).Decode(&result)

	if err != nil {
		log.Println(err)
	}

	return result
}

func (r tokenRepository) GetAll() []*domain.Token {
	cur, _ := r.getCollection().Find(context.Background(), bson.D{{}})
	var results []*domain.Token

	for cur.Next(context.TODO()) {
		var elem domain.Token
		err := cur.Decode(&elem)
		if err != nil {
			log.Println(err)
		}
		results = append(results, &elem)
	}

	return results
}

func (r tokenRepository) getCollection() *mongo.Collection {
	return r.client.Database(os.Getenv("MONGODB_DATABASE")).Collection("tokens")
}
