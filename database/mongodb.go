package database

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDBClient() *mongo.Client {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "mongo"
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "mongo"
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "mongo"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "mongo"
	}

	connectionUrl := fmt.Sprintf("mongodb+srv://%s:%s@%s/?retryWrites=true&w=majority", dbUser, dbPassword, dbHost)
	opts := options.Client().ApplyURI(connectionUrl).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		panic(err)
	}

	if err := client.Database(dbName).RunCommand(context.Background(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		panic(err)
	}
	return client
}
