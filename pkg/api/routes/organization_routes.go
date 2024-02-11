package routes

import (
	"ideanest/pkg/api/handlers"
	"ideanest/pkg/api/middleware"

	"github.com/gin-gonic/gin"
)

func OrganizationRoute(router *gin.Engine) {
	organizationGroup := router.Group("/organization")
	organizationGroup.Use(middleware.Auth())
	organizationGroup.POST("/", handlers.CreateOrganization())
	organizationGroup.GET("/", handlers.ReadOrganizations())

	idGroup := organizationGroup.Group("/:id")
	idGroup.Use(middleware.AccessAuth())
	idGroup.GET("/", handlers.ReadOrganization())
	idGroup.PUT("/", handlers.UpdateOrganization())
	idGroup.DELETE("/", handlers.DeleteOrganization())
	idGroup.POST("/invite", handlers.InviteUser())
}
