package main

import (
	"context"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	MongoDBConnectionStringKey = "MONGODB_URI"
)

func main() {
	mongoURI := os.Getenv(MongoDBConnectionStringKey)
	mdb, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		panic(err.Error())
	}
	if err := mdb.Ping(context.Background(), nil); err != nil {
		panic(err.Error())
	}
	defer mdb.Disconnect(context.Background())
}
