package middleware

import (
	"ideanest/pkg/controllers"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("Authorization")
		if clientToken == "" {
			c.JSON(
				http.StatusForbidden,
				gin.H{"message": "No Authorization header provided"})
			c.Abort()
			return
		}

		extractedToken := strings.Split(clientToken, "Bearer ")
		if len(extractedToken) != 2 {
			// If the token is not in the correct format, return a 400 status code
			c.JSON(
				http.StatusBadRequest,
				gin.H{"message": "Incorrect Format of Authorization Token"})
			c.Abort()
			return
		}

		// Trim the token
		clientToken = strings.TrimSpace(extractedToken[1])

		// validate token
		email, err := controllers.ValidateToken(clientToken)
		if err != nil {
			c.JSON(
				http.StatusUnauthorized,
				gin.H{"message": "Token is not authorized"})
			c.Abort()
			return
		}

		// get the user
		user, err := controllers.RetreiveUserByEmail(email)
		if err != nil {
			c.JSON(
				http.StatusUnauthorized,
				gin.H{"message": "Token is not bound to a current user"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
