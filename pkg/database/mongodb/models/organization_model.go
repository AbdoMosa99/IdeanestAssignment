package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Model representing an Organization Entity
type OrganizationModel struct {
	Id                  primitive.ObjectID `json:"organization_id"`
	Name                string             `json:"name"`
	Description         string             `json:"description"`
	OrganizationMembers []UserReference    `json:"organization_members"`
}

// Model representing the organization member user
type UserReference struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	AcessLevel string `json:"access_level"`
}
