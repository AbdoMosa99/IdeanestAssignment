package middleware

import (
	"ideanest/pkg/controllers"
	"ideanest/pkg/database/mongodb/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AccessAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		res, _ := ctx.Get("user")
		user := res.(*models.UserModel)

		id := ctx.Param("id")
		organization, err := controllers.RetrieveOrganization(id)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "Invalid ID"})
			ctx.Abort()
			return
		}

		accessLevel := controllers.GetUserAccess(user, organization)
		if accessLevel == nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Not allowed"})
			ctx.Abort()
			return
		}

		if *accessLevel == "readonly" && ctx.Request.Method != http.MethodGet {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Not allowed"})
			ctx.Abort()
			return
		}

		ctx.Set("organization", organization)
		ctx.Next()
	}
}
