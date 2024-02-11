package routes

import (
	"ideanest/pkg/api/handlers"
	"ideanest/pkg/api/middleware"

	"github.com/gin-gonic/gin"
)

func OrganizationRoute(router *gin.Engine) {
	group := router.Group("/organization").Use(middleware.Auth())
	group.POST("/", handlers.CreateOrganization())
	group.GET("/", handlers.ReadOrganizations())
	group.GET("/:id", handlers.ReadOrganization())
	group.PUT("/:id", handlers.UpdateOrganization())
	group.DELETE("/:id", handlers.DeleteOrganization())
}
