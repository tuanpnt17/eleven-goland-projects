package database

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
	"log"
	"os"
	"time"
)

func dbInstance() *mongo.Client {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	mongoDbConnectionString := os.Getenv("MONGODB_URL") // Todo: Add MONGODB_URL to env

	client, err := mongo.Connect(options.Client().ApplyURI(mongoDbConnectionString))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	return client
}

var Client *mongo.Client = dbInstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	databaseName := os.Getenv("MONGODB_DATABASE") // Todo: Add MONGODB_DATABASE to .env
	return client.Database(databaseName).Collection(collectionName)
}

func CloseConnection() {
	err := Client.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
