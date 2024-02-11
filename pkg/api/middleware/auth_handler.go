package middleware

import (
	"ideanest/pkg/controllers"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Authentication Handler
func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get the authorization from header
		clientToken := ctx.Request.Header.Get("Authorization")
		if clientToken == "" {
			// if no token was provided
			ctx.JSON(
				http.StatusForbidden,
				gin.H{"message": "No Authorization header provided"})
			ctx.Abort()
			return
		}

		// extract the token from header (format: "Bearer {token}")
		extractedToken := strings.Split(clientToken, "Bearer ")
		if len(extractedToken) != 2 {
			// If the token is not in the correct format, return a 400 status code
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{"message": "Incorrect Format of Authorization Token"})
			ctx.Abort()
			return
		}

		// Trim the token
		clientToken = strings.TrimSpace(extractedToken[1])

		// validate token
		email, err := controllers.ValidateToken(clientToken)
		if err != nil {
			ctx.JSON(
				http.StatusUnauthorized,
				gin.H{"message": "Token is not authorized"})
			ctx.Abort()
			return
		}

		// get the user attached to this token
		user, err := controllers.RetreiveUserByEmail(email)
		if err != nil {
			ctx.JSON(
				http.StatusUnauthorized,
				gin.H{"message": "Token is not bound to a current user"})
			ctx.Abort()
			return
		}

		// set user in context for next handler
		ctx.Set("user", user)

		// forward to next handler
		ctx.Next()
	}
}
