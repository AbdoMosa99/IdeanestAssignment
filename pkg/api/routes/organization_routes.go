package routes

import (
	"ideanest/pkg/controllers"

	"github.com/gin-gonic/gin"
)

func OrganizationRoute(router *gin.Engine) {
	router.POST("/organization", controllers.CreateOrganization)
	router.GET("/organization", controllers.ReadOrganizations)
	router.GET("/organization/:id", controllers.ReadOrganization)
	router.PUT("/organization/:id", controllers.UpdateOrganization)
	router.DELETE("/organization/:id", controllers.DeleteOrganization)
}
