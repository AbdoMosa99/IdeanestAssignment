package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
)

var db *mongo.Client
var UserCollection *mongo.Collection
var OrganizationCollection *mongo.Collection

func getCollection(collectionName string) *mongo.Collection {
	collection := db.Database("ideanest").Collection(collectionName)
	return collection
}
