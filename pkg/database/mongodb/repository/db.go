package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// database state (maily collections)
var db *mongo.Client
var UserCollection *mongo.Collection
var OrganizationCollection *mongo.Collection

// get the database collection with the provided name
func getCollection(collectionName string) *mongo.Collection {
	// TODO: should be read from config
	dbName := "ideanest"

	collection := db.Database(dbName).Collection(collectionName)
	return collection
}
