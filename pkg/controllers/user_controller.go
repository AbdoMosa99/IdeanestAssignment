package controllers

import (
	"context"
	"ideanest/pkg/database/mongodb/models"
	"ideanest/pkg/database/mongodb/repository"

	"go.mongodb.org/mongo-driver/bson"
)

func InsertUser(user models.UserModel) error {
	_, err := repository.UserCollection.InsertOne(context.TODO(), user)
	return err
}

func RetreiveUserByEmail(email string) (*models.UserModel, error) {
	var user models.UserModel
	err := repository.UserCollection.FindOne(
		context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
