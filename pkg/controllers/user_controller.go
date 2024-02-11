package controllers

import (
	"context"
	"ideanest/pkg/database/mongodb/models"
	"ideanest/pkg/database/mongodb/repository"

	"go.mongodb.org/mongo-driver/bson"
)

// Insert a given user into the database
func InsertUser(user models.UserModel) error {
	_, err := repository.UserCollection.InsertOne(context.TODO(), user)
	return err
}

// Get the user object with the given email
func RetreiveUserByEmail(email string) (*models.UserModel, error) {
	var user models.UserModel

	// find user
	filter := bson.M{"email": email}
	err := repository.UserCollection.FindOne(context.TODO(), filter).Decode(&user)

	// email not found
	if err != nil {
		return nil, err
	}

	return &user, nil
}
