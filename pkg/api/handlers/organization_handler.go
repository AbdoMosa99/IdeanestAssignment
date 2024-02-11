package handlers

import (
	"ideanest/pkg/controllers"
	"ideanest/pkg/database/mongodb/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateOrganization() gin.HandlerFunc {
	type bodyPayload struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description" binding:"required"`
	}
	return func(ctx *gin.Context) {
		// check that request body matches our model
		var bodyData bodyPayload
		err := ctx.BindJSON(&bodyData)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Missing required body fields"})
			ctx.Abort()
			return
		}

		// if valid, create it
		organization := models.OrganizationModel{
			Name:        bodyData.Name,
			Description: bodyData.Description,
		}
		err = controllers.InsertOrganization(&organization)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
			ctx.Abort()
			return
		}
		ctx.JSON(http.StatusCreated, gin.H{"organization_id": organization.Id})
	}
}

func ReadOrganizations() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		organizations, err := controllers.ListOrganizations()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
			ctx.Abort()
			return
		}
		ctx.JSON(http.StatusOK, organizations)
	}
}

func ReadOrganization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		organization, err := controllers.RetrieveOrganization(id)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "Invalid ID"})
			ctx.Abort()
			return
		}
		ctx.JSON(http.StatusOK, organization)
	}
}

func UpdateOrganization() gin.HandlerFunc {
	type bodyPayload struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description" binding:"required"`
	}

	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		// check that request body matches our model
		var bodyData bodyPayload
		err := ctx.BindJSON(&bodyData)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Missing required body fields"})
			ctx.Abort()
			return
		}

		organization := models.OrganizationModel{
			Name:        bodyData.Name,
			Description: bodyData.Description,
		}
		err = controllers.EditOrganization(id, &organization)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "Invalid ID"})
			ctx.Abort()
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
		id := ctx.Param("id")
		err := controllers.RemoveOrganization(id)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "Invalid ID"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "deleted"})
	}
}

func InviteUser() gin.HandlerFunc {
	type bodyPayload struct {
		UserEmail string `json:"user_email" binding:"required"`
	}

	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		var bodyData bodyPayload
		err := ctx.BindJSON(&bodyData)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Missing required body fields"})
			ctx.Abort()
			return
		}

		user, err := controllers.RetreiveUserByEmail(bodyData.UserEmail)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "No user found with given user_email"})
			ctx.Abort()
			return
		}

		organization, err := controllers.RetrieveOrganization(id)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "Invalid ID"})
			ctx.Abort()
			return
		}

		added, err := controllers.AddUserToOrganization(user, organization)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
			ctx.Abort()
			return
		}

		if added {
			ctx.JSON(http.StatusOK, gin.H{"message": "Invited successfully"})
		} else {
			ctx.JSON(http.StatusOK, gin.H{"message": "Already invited"})
		}
	}
}
