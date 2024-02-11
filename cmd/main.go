package main

import (
	"context"
	"ideanest/pkg/api/routes"
	"ideanest/pkg/database/mongodb/repository"

	"github.com/gin-gonic/gin"
)

// entrypoint
func main() {
	// initialize database
	ctx := context.Background()
	repository.ConnectDB(ctx)
	defer repository.DisconnectDB()

	// set up router
	router := gin.Default()
	routes.UserRoute(router)
	routes.OrganizationRoute(router)

	// run on port 8080
	router.Run(":8080")
}
