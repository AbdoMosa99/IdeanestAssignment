package handlers

import (
	"ideanest/pkg/controllers"
	"ideanest/pkg/database/mongodb/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateOrganization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// check that request body matches our model
		var organization models.OrganizationModel
		err := ctx.BindJSON(&organization)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			ctx.Abort()
			return
		}

		// if valid, create it
		err = controllers.InsertOrganization(&organization)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
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
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
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
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
			ctx.Abort()
			return
		}
		ctx.JSON(http.StatusOK, organization)
	}
}

func UpdateOrganization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		// check that request body matches our model
		var organization models.OrganizationModel
		err := ctx.BindJSON(&organization)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			ctx.Abort()
			return
		}

		err = controllers.EditOrganization(id, &organization)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
			ctx.Abort()
			return
		}

		ctx.JSON(http.StatusOK, organization)
	}
}

func DeleteOrganization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		err := controllers.RemoveOrganization(id)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "deleted"})
	}
}
