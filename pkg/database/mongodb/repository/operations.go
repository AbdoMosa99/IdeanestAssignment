package repository

import (
	"ideanest/pkg/database/mongodb/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateOrganization(organization *models.OrganizationModel) error {
	organization.Id = primitive.NewObjectID()
	_, err := organizationCollection.InsertOne(ctx, organization)
	return err
}
