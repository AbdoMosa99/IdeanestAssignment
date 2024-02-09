package routes

import (
	"ideanest/pkg/controllers"

	"github.com/gin-gonic/gin"
)

func OrganizationRoute(router *gin.Engine) {
	router.POST("/organizations", controllers.CreateOrganization)
}
