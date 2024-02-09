package repository

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func ConnectDB(c context.Context) *mongo.Client {
	client, err := mongo.Connect(c,
		options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	// test database connection
	err = client.Ping(c, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	// initialize state
	ctx = c
	db = client
	organizationCollection = GetCollection("organizations")

	return db
}

func DisconnectDB() {
	err := db.Disconnect(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func GetCollection(collectionName string) *mongo.Collection {
	collection := db.Database("ideanest").Collection(collectionName)
	return collection
}
