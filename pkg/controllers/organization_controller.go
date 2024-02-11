package controllers

import (
	"context"
	"ideanest/pkg/database/mongodb/models"
	"ideanest/pkg/database/mongodb/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InsertOrganization(organization *models.OrganizationModel) error {
	organization.Id = primitive.NewObjectID()
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
	update := bson.M{"name": organization.Name, "description": organization.Description}
	_, err = repository.OrganizationCollection.UpdateOne(
		context.TODO(), bson.M{"id": objID}, bson.M{"$set": update})
	return err
}

func RemoveOrganization(id string) error {
	objID, _ := primitive.ObjectIDFromHex(id)
	_, err := repository.OrganizationCollection.DeleteOne(
		context.TODO(), bson.M{"id": objID})
	return err
}
