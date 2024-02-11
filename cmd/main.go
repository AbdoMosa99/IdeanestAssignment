package main

import (
	"context"
	"ideanest/pkg/api/routes"
	"ideanest/pkg/database/mongodb/repository"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()
	repository.ConnectDB(ctx)
	defer repository.DisconnectDB()

	router := gin.Default()
	routes.OrganizationRoute(router)
	routes.UserRoute(router)
	router.Run(":8080")
}
