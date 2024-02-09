package controllers

import (
	"ideanest/pkg/database/mongodb/models"
	"ideanest/pkg/database/mongodb/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateOrganization(ctx *gin.Context) {
	// check that request body matches our model
	var organization models.OrganizationModel
	err := ctx.BindJSON(&organization)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// if valid, create it
	err = repository.InsertOrganization(&organization)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"organization_id": organization.Id})
}

func ReadOrganizations(ctx *gin.Context) {
	organizations, err := repository.ListOrganizations()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, organizations)
}

func ReadOrganization(ctx *gin.Context) {
	id := ctx.Param("id")
	organization, err := repository.RetrieveOrganization(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}
	ctx.JSON(http.StatusOK, organization)
}

func UpdateOrganization(ctx *gin.Context) {
	id := ctx.Param("id")

	// check that request body matches our model
	var organization models.OrganizationModel
	err := ctx.BindJSON(&organization)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = repository.EditOrganization(id, &organization)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}

	ctx.JSON(http.StatusOK, organization)
}

func DeleteOrganization(ctx *gin.Context) {
	id := ctx.Param("id")
	err := repository.RemoveOrganization(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
