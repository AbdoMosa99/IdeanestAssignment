package handlers

import (
	"ideanest/pkg/controllers"
	"ideanest/pkg/database/mongodb/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler for creating an organization
func CreateOrganization() gin.HandlerFunc {
	// body format
	type bodyInterface struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description" binding:"required"`
	}

	return func(ctx *gin.Context) {
		// check that request body matches our format
		var bodyData bodyInterface
		err := ctx.BindJSON(&bodyData)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Missing required body fields"})
			return
		}

		// get current user
		res, _ := ctx.Get("user")
		user := res.(*models.UserModel)

		// create the organization with given data
		organization := models.OrganizationModel{
			Name:        bodyData.Name,
			Description: bodyData.Description,
			OrganizationMembers: []models.UserReference{
				{
					Name:       user.Name,
					Email:      user.Email,
					AcessLevel: "owner",
				},
			},
		}

		// insert it into the database
		err = controllers.InsertOrganization(&organization)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
			return
		}

		// return success
		ctx.JSON(http.StatusCreated, gin.H{"organization_id": organization.Id})
	}
}

// Get all organizations (that user can read)
func ReadOrganizations() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// get current user
		res, _ := ctx.Get("user")
		user := res.(*models.UserModel)

		// get all organizations from database
		allOrganizations, err := controllers.ListOrganizations()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
			return
		}

		// filter them according the user access
		userOrganizations := []models.OrganizationModel{}
		for _, organization := range allOrganizations {
			if controllers.GetUserAccess(user, &organization) != nil {
				userOrganizations = append(userOrganizations, organization)
			}
		}

		// return success
		ctx.JSON(http.StatusOK, userOrganizations)
	}
}

// Read an organization with the given ID
func ReadOrganization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// from middleware, we know that the user is granted
		res, _ := ctx.Get("organization")
		organization := res.(*models.OrganizationModel)

		// return success
		ctx.JSON(http.StatusOK, organization)
	}
}

// Update an organization with a given ID
func UpdateOrganization() gin.HandlerFunc {
	// body format
	type bodyInterface struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description" binding:"required"`
	}

	return func(ctx *gin.Context) {
		// check that request body matches our format
		var bodyData bodyInterface
		err := ctx.BindJSON(&bodyData)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Missing required body fields"})
			return
		}

		// from middleware, we know user is granted
		res, _ := ctx.Get("organization")
		organization := res.(*models.OrganizationModel)

		// update values to those in the body
		organization.Name = bodyData.Name
		organization.Description = bodyData.Description

		// update them into the database
		err = controllers.EditOrganization(organization)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
			return
		}

		// return success
		ctx.JSON(http.StatusOK, gin.H{
			"organization_id": organization.Id,
			"name":            organization.Name,
			"description":     organization.Description,
		})
	}
}

// delete an organization
func DeleteOrganization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// from middleware, we know user is granted
		res, _ := ctx.Get("organization")
		organization := res.(*models.OrganizationModel)

		// remove it from database
		err := controllers.RemoveOrganization(organization)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Something went Wrong"})
			return
		}

		// return success
		ctx.JSON(http.StatusOK, gin.H{"message": "deleted"})
	}
}

// invite user to be able to read a given organization
func InviteUser() gin.HandlerFunc {
	// body format
	type bodyInterface struct {
		UserEmail string `json:"user_email" binding:"required"`
	}

	return func(ctx *gin.Context) {
		// check that body matches our format
		var bodyData bodyInterface
		err := ctx.BindJSON(&bodyData)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Missing required body fields"})
			return
		}

		// from middleware, we know user is granted
		res, _ := ctx.Get("organization")
		organization := res.(*models.OrganizationModel)

		// get the invited user
		user, err := controllers.RetreiveUserByEmail(bodyData.UserEmail)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "No user found with given user_email"})
			return
		}

		// add the invited user to the organization members list
		added, err := controllers.AddUserToOrganization(user, organization)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
			return
		}

		// return success
		if added {
			ctx.JSON(http.StatusOK, gin.H{"message": "Invited successfully"})
		} else {
			ctx.JSON(http.StatusOK, gin.H{"message": "Already invited"})
		}
	}
}
