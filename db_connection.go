package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

var client *mongo.Client = nil

func GetClient() *mongo.Client {

	if client == nil {

		url := os.Getenv("MONGODB_URL")

		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

		newClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+url))

		if err != nil {
			log.Fatal("Failed To Connect to mongo")
		} else {
			log.Printf("Connected to mongo")
		}

		client = newClient
	}

	return client
}
