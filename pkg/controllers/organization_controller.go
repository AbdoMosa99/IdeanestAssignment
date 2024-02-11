package controllers

import (
	"context"
	"ideanest/pkg/database/mongodb/models"
	"ideanest/pkg/database/mongodb/repository"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Insert a given organization into the mongodb database
func InsertOrganization(organization *models.OrganizationModel) error {
	// get a new object id
	organization.Id = primitive.NewObjectID()

	_, err := repository.OrganizationCollection.InsertOne(context.TODO(), organization)
	return err
}

// Get and return all organization objects/documents from database
func ListOrganizations() ([]models.OrganizationModel, error) {
	// get a new context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// get a cursor with no filters
	cur, err := repository.OrganizationCollection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx) // make sure we close that cursor

	// get all organizations
	var organizations []models.OrganizationModel
	err = cur.All(ctx, &organizations)

	// if no organizations found, make it an empty list (not nil)
	if organizations == nil {
		organizations = []models.OrganizationModel{}
	}

	return organizations, err
}

// Get an organization object from the database with a given id
func RetrieveOrganization(id string) (*models.OrganizationModel, error) {
	var organization models.OrganizationModel

	// convert id from string to ObjectID
	objID, _ := primitive.ObjectIDFromHex(id)

	// find the document with that id and convert it to organization object
	err := repository.OrganizationCollection.FindOne(context.TODO(),
		bson.M{"id": objID}).Decode(&organization)
	if err != nil {
		return nil, err
	}

	return &organization, nil
}

// Update a given organization in database with the specified key to the given data
func EditOrganization(organization *models.OrganizationModel) error {
	filter := bson.M{"id": organization.Id}
	update := bson.M{"$set": organization}

	_, err := repository.OrganizationCollection.UpdateOne(context.TODO(), filter, update)

	return err
}

// Delete an organization from database
func RemoveOrganization(organization *models.OrganizationModel) error {
	filter := bson.M{"id": organization.Id}

	_, err := repository.OrganizationCollection.DeleteOne(context.TODO(), filter)

	return err
}

// Insert a user as a readonly user to the given organization
// Returns true if added and false if already existed
func AddUserToOrganization(user *models.UserModel, organization *models.OrganizationModel) (bool, error) {
	// check if the user is already existed
	for _, existedUser := range organization.OrganizationMembers {
		if existedUser.Email == user.Email {
			return false, nil
		}
	}

	// create and append a reference to that user
	userRef := models.UserReference{
		Name:       user.Name,
		Email:      user.Email,
		AcessLevel: "readonly",
	}
	organization.OrganizationMembers = append(organization.OrganizationMembers, userRef)

	// update the data into the database
	filter := bson.M{"id": organization.Id}
	update := bson.M{"$set": organization}
	_, err := repository.OrganizationCollection.UpdateOne(context.TODO(), filter, update)

	return true, err
}

// Get the user access level of a given organization (nil if no access)
func GetUserAccess(user *models.UserModel, organization *models.OrganizationModel) *string {
	// Search for the user
	for _, userRef := range organization.OrganizationMembers {
		// if existed, return its access level
		if user.Email == userRef.Email {
			return &userRef.AcessLevel
		}
	}

	// if not found,
	return nil
}
