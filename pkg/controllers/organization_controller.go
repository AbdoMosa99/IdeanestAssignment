package controllers

import (
	"context"
	"ideanest/pkg/database/mongodb/models"
	"ideanest/pkg/database/mongodb/repository"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InsertOrganization(organization *models.OrganizationModel) error {
	organization.Id = primitive.NewObjectID()
	_, err := repository.OrganizationCollection.InsertOne(context.TODO(), organization)
	return err
}

func ListOrganizations() ([]models.OrganizationModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cur, err := repository.OrganizationCollection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var organizations []models.OrganizationModel
	err = cur.All(ctx, &organizations)
	if organizations == nil {
		organizations = []models.OrganizationModel{}
	}

	return organizations, err
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

func EditOrganization(organization *models.OrganizationModel) error {
	_, err := repository.OrganizationCollection.UpdateOne(
		context.TODO(), bson.M{"id": organization.Id}, bson.M{"$set": organization})
	return err
}

func RemoveOrganization(organization *models.OrganizationModel) error {
	_, err := repository.OrganizationCollection.DeleteOne(
		context.TODO(), bson.M{"id": organization.Id})
	return err
}

func AddUserToOrganization(user *models.UserModel, organization *models.OrganizationModel) (bool, error) {
	for _, existedUser := range organization.OrganizationMembers {
		if existedUser.Email == user.Email {
			return false, nil
		}
	}
	userRef := models.UserReference{
		Name:       user.Name,
		Email:      user.Email,
		AcessLevel: "readonly",
	}
	organization.OrganizationMembers = append(organization.OrganizationMembers, userRef)
	_, err := repository.OrganizationCollection.UpdateOne(
		context.TODO(), bson.M{"id": organization.Id}, bson.M{"$set": organization})
	return true, err
}

func GetUserAccess(user *models.UserModel, organization *models.OrganizationModel) *string {
	for _, userRef := range organization.OrganizationMembers {
		if user.Email == userRef.Email {
			return &userRef.AcessLevel
		}
	}
	return nil
}
