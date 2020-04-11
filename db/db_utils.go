package db

import (
	"context"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

var _ = godotenv.Load()

var client = getClient()

func getClient() *mongo.Client {

	url := os.Getenv("MONGODB_URL")

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://"+url))

	if err != nil {
		log.Printf("Failed To Connect to mongo %s", err)
	} else {
		log.Println("Connected to mongo")
	}

	return client
}

func getCollection() *mongo.Collection {
	return client.Database(os.Getenv("MONGODB_DATABASE")).Collection("tokens")
}
