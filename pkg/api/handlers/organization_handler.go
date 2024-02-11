package handlers

import (
	"ideanest/pkg/controllers"
	"ideanest/pkg/database/mongodb/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateOrganization() gin.HandlerFunc {
	type bodyInterface struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description" binding:"required"`
	}
	return func(ctx *gin.Context) {
		// check that request body matches our model
		var bodyData bodyInterface
		err := ctx.BindJSON(&bodyData)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Missing required body fields"})
			return
		}

		res, _ := ctx.Get("user")
		user := res.(*models.UserModel)

		// if valid, create it
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
		err = controllers.InsertOrganization(&organization)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
			return
		}
		ctx.JSON(http.StatusCreated, gin.H{"organization_id": organization.Id})
	}
}

func ReadOrganizations() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		res, exists := ctx.Get("user")
		if !exists {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
			return
		}
		user := res.(*models.UserModel)

		allOrganizations, err := controllers.ListOrganizations()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
			return
		}

		userOrganizations := []models.OrganizationModel{}
		for _, organization := range allOrganizations {
			if controllers.GetUserAccess(user, &organization) != nil {
				userOrganizations = append(userOrganizations, organization)
			}
		}

		ctx.JSON(http.StatusOK, userOrganizations)
	}
}

func ReadOrganization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		res, _ := ctx.Get("organization")
		organization := res.(*models.OrganizationModel)
		ctx.JSON(http.StatusOK, organization)
	}
}

func UpdateOrganization() gin.HandlerFunc {
	type bodyInterface struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description" binding:"required"`
	}

	return func(ctx *gin.Context) {
		// check that request body matches our model
		var bodyData bodyInterface
		err := ctx.BindJSON(&bodyData)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Missing required body fields"})
			return
		}

		res, _ := ctx.Get("organization")
		organization := res.(*models.OrganizationModel)

		organization.Name = bodyData.Name
		organization.Description = bodyData.Description
		err = controllers.EditOrganization(organization)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"organization_id": organization.Id,
			"name":            organization.Name,
			"description":     organization.Description,
		})
	}
}

func DeleteOrganization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		res, _ := ctx.Get("organization")
		organization := res.(*models.OrganizationModel)
		err := controllers.RemoveOrganization(organization)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Something went Wrong"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "deleted"})
	}
}

func InviteUser() gin.HandlerFunc {
	type bodyInterface struct {
		UserEmail string `json:"user_email" binding:"required"`
	}

	return func(ctx *gin.Context) {
		var bodyData bodyInterface
		err := ctx.BindJSON(&bodyData)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Missing required body fields"})
			return
		}

		res, _ := ctx.Get("organization")
		organization := res.(*models.OrganizationModel)

		user, err := controllers.RetreiveUserByEmail(bodyData.UserEmail)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "No user found with given user_email"})
			return
		}

		added, err := controllers.AddUserToOrganization(user, organization)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
			return
		}

		if added {
			ctx.JSON(http.StatusOK, gin.H{"message": "Invited successfully"})
		} else {
			ctx.JSON(http.StatusOK, gin.H{"message": "Already invited"})
		}
	}
}
