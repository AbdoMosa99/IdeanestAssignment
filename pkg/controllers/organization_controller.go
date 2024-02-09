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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// if valid, create it
	err = repository.CreateOrganization(&organization)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, organization)
}
