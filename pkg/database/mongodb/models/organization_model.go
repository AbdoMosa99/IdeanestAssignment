package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type OrganizationModel struct {
	Id                  primitive.ObjectID `json:"organization_id"`
	Name                string             `json:"name"`
	Description         string             `json:"description"`
	OrganizationMembers []UserModel        `json:"organization_members"`
}
