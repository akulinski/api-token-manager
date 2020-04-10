package db

import (
	"context"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

var _ = godotenv.Load()


type Token struct {
	Id       string    `json:"id" bson:"_id,omitempty"`
	IssuedAt time.Time `json:"issuedAt" bson:"issuedAt,omitempty"`
	Issuer   string    `json:"issuer" bson:"issuer,omitempty"`
	UserID   string    `json:"userID" bson:"userID"`
	Token    string    `json:"token" bson:"token"`
	Expired  bool      `json:"expired" bson:"expired,omitempty"`
	Revoked  bool      `json:"revoked" bson:"revoked,omitempty"`
}

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

func Insert(token Token) *mongo.InsertOneResult {
	collection := getCollection()
	res, err := collection.InsertOne(context.Background(), token)

	if err != nil {
		log.Printf("Failed to insert token %s", err)
	}

	return res
}

func GetAll() []*Token {
	cur, _ := getCollection().Find(context.Background(), bson.D{{}})
	var results []*Token

	for cur.Next(context.TODO()) {
		var elem Token
		err := cur.Decode(&elem)
		if err != nil {
			log.Println(err)
		}
		results = append(results, &elem)
	}

	return results
}

func GetById(id primitive.ObjectID) Token {

	filter := bson.M{"_id": id}

	var result Token

	err := getCollection().FindOne(context.Background(), filter).Decode(&result)

	if err != nil {
		log.Println(err)
	}

	return result
}

func RevokeToken(id primitive.ObjectID) int64 {

	filter := bson.M{"_id": id}

	byId := GetById(id)

	byId.Revoked = true

	result, err := getCollection().UpdateOne(context.Background(), filter, bson.D{
		{
			"$set", bson.D{{"revoked", true}},
		},
	})

	if err != nil {
		log.Println(err)
	}

	return result.ModifiedCount
}

func getCollection() *mongo.Collection {
	return client.Database(os.Getenv("MONGODB_DATABASE")).Collection("tokens")
}
