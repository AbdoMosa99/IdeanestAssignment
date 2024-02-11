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
	db = client
	OrganizationCollection = getCollection("organizations")
	UserCollection = getCollection("users")

	return db
}

func DisconnectDB() {
	err := db.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
}
