package middleware

import (
	"ideanest/pkg/controllers"
	"ideanest/pkg/database/mongodb/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Authorization Handler
// Middleware handler which checks that the user has access to the asked resource
func AccessAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// get the user (Auth handler already checked authentication)
		res, _ := ctx.Get("user")
		user := res.(*models.UserModel)

		// get resource and check it exists
		id := ctx.Param("id")
		organization, err := controllers.RetrieveOrganization(id)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "Invalid ID"})
			ctx.Abort()
			return
		}

		// get the access level of current user on that resource
		accessLevel := controllers.GetUserAccess(user, organization)

		// if user has no access at all
		if accessLevel == nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Not allowed"})
			ctx.Abort()
			return
		}

		// if user has readonly access and tries to do write operation
		if *accessLevel == "readonly" && ctx.Request.Method != http.MethodGet {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Not allowed"})
			ctx.Abort()
			return
		}

		// if granted access, set the resource in the context
		ctx.Set("organization", organization)

		// forward to the next handler
		ctx.Next()
	}
}
