package handlers

import (
	"ideanest/pkg/controllers"
	"ideanest/pkg/database/mongodb/models"
	"ideanest/pkg/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SignUp() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user models.UserModel

		err := ctx.BindJSON(&user)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Request body is invalid"})
			ctx.Abort()
			return
		}

		user.Password = utils.HashPassword(user.Password)
		err = controllers.InsertUser(user)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			ctx.Abort()
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Registered Successfully"})
	}
}

func SignIn() gin.HandlerFunc {
	type loginPayload struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	return func(ctx *gin.Context) {
		var loginData loginPayload

		err := ctx.BindJSON(&loginData)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			ctx.Abort()
			return
		}

		user, err := controllers.RetreiveUser(loginData.Email)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid credentials"})
			ctx.Abort()
			return
		}

		err = utils.CheckPassword(loginData.Password, user.Password)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid credentials"})
			ctx.Abort()
			return
		}

		access_token, err := controllers.GenerateToken(user.Email, 24*time.Hour)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			ctx.Abort()
			return
		}
		refresh_token, err := controllers.GenerateToken(user.Email, 30*24*time.Hour)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			ctx.Abort()
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message":       "Logged in successfully",
			"access_token":  access_token,
			"refresh_token": refresh_token,
		})

	}
}

func RefreshToken() gin.HandlerFunc {
	type refreshPayload struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	return func(ctx *gin.Context) {
		var refreshData refreshPayload

		err := ctx.BindJSON(&refreshData)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			ctx.Abort()
			return
		}

		refreshToken := refreshData.RefreshToken
		email, err := controllers.ValidateToken(refreshToken)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			ctx.Abort()
			return
		}

		user, err := controllers.RetreiveUser(email)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user"})
			ctx.Abort()
			return
		}

		accessToken, err := controllers.GenerateToken(user.Email, 24*time.Hour)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			ctx.Abort()
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message":       "Refreshed successfully",
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})
	}
}
