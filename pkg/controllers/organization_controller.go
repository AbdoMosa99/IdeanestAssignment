package controllers

import (
	"context"
	"ideanest/pkg/database/mongodb/models"
	"ideanest/pkg/database/mongodb/repository"

	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InsertOrganization(organization *models.OrganizationModel) error {
	organization.Id = primitive.NewObjectID()
	organization.OrganizationMembers = []models.UserModel{}
	_, err := repository.OrganizationCollection.InsertOne(context.TODO(), organization)
	return err
}

func ListOrganizations() ([]models.OrganizationModel, error) {
	cur, err := repository.OrganizationCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}

	var organizations []models.OrganizationModel
	err = cur.All(context.TODO(), &organizations)
	if err != nil {
		return nil, err
	}

	return organizations, nil
}

func RetrieveOrganization(id string) (*models.OrganizationModel, error) {
	var organization models.OrganizationModel
	objID, _ := primitive.ObjectIDFromHex(id)
	err := repository.OrganizationCollection.FindOne(context.TODO(),
		bson.M{"id": objID}).Decode(&organization)
	if err != nil {
		return nil, err
	}
	return &organization, nil
}

func EditOrganization(id string, organization *models.OrganizationModel) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	organization.Id = objID
	_, err = repository.OrganizationCollection.UpdateOne(
		context.TODO(), bson.M{"id": objID}, bson.M{"$set": organization})
	return err
}

func RemoveOrganization(id string) error {
	objID, _ := primitive.ObjectIDFromHex(id)
	res, err := repository.OrganizationCollection.DeleteOne(
		context.TODO(), bson.M{"id": objID})
	if err != nil {
		return err
	}
	if res.DeletedCount != 1 {
		return errors.New("resource not found")
	}
	return nil
}

func AddUserToOrganization(user *models.UserModel, organization *models.OrganizationModel) (bool, error) {
	for _, existedUser := range organization.OrganizationMembers {
		if existedUser.Email == user.Email {
			return false, nil
		}
	}
	organization.OrganizationMembers = append(organization.OrganizationMembers, *user)
	_, err := repository.OrganizationCollection.UpdateOne(
		context.TODO(), bson.M{"id": organization.Id}, bson.M{"$set": organization})
	return true, err
}
