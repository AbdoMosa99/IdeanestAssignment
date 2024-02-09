package repository

import (
	"ideanest/pkg/database/mongodb/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InsertOrganization(organization *models.OrganizationModel) error {
	organization.Id = primitive.NewObjectID()
	_, err := organizationCollection.InsertOne(ctx, organization)
	return err
}

func ListOrganizations() ([]models.OrganizationModel, error) {
	cur, err := organizationCollection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var organizations []models.OrganizationModel
	err = cur.All(ctx, &organizations)
	if err != nil {
		return nil, err
	}

	return organizations, nil
}

func RetrieveOrganization(id string) (*models.OrganizationModel, error) {
	var organization models.OrganizationModel
	objID, _ := primitive.ObjectIDFromHex(id)
	err := organizationCollection.FindOne(ctx, bson.M{"id": objID}).Decode(&organization)
	if err != nil {
		return nil, err
	}
	return &organization, nil
}

func EditOrganization(id string, organization *models.OrganizationModel) error {
	objID, _ := primitive.ObjectIDFromHex(id)
	organization.Id = objID
	update := bson.M{"name": organization.Name, "description": organization.Description}
	_, err := organizationCollection.UpdateOne(ctx, bson.M{"id": objID}, bson.M{"$set": update})
	return err
}

func RemoveOrganization(id string) error {
	objID, _ := primitive.ObjectIDFromHex(id)
	_, err := organizationCollection.DeleteOne(ctx, bson.M{"id": objID})
	return err
}
