package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type OrganizationModel struct {
	Id                  primitive.ObjectID `json:"organization_id"`
	Name                string             `json:"name"`
	Description         string             `json:"description"`
	OrganizationMembers []UserReference    `json:"organization_members"`
}

type UserReference struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	AcessLevel string `json:"access_level"`
}
