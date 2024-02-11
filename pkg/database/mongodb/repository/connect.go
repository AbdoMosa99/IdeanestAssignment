package repository

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Function that connects to the mongodb database and initializes the state
// should be called in the starting of the application
func ConnectDB(c context.Context) *mongo.Client {
	// TODO: should be loaded from config
	dbURI := "mongodb://localhost:27017"

	client, err := mongo.Connect(c, options.Client().ApplyURI(dbURI))
	if err != nil {
		// can't connect!
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

// function that disconnects from the database
// should be called when the application stops
func DisconnectDB() {
	err := db.Disconnect(context.TODO())
	if err != nil {
		// can't disconnect
		log.Fatal(err)
	}
}
