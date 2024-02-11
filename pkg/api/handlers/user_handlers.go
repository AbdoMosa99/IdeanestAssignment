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
	type bodyPayload struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	return func(ctx *gin.Context) {
		var bodyData bodyPayload
		err := ctx.BindJSON(&bodyData)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Missing required body fields"})
			ctx.Abort()
			return
		}

		user := models.UserModel{
			Name:        bodyData.Name,
			Email:       bodyData.Email,
			AccessLevel: "normal",
			Password:    utils.HashPassword(bodyData.Password),
		}
		err = controllers.InsertUser(user)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
			ctx.Abort()
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Registered Successfully"})
	}
}

func SignIn() gin.HandlerFunc {
	type bodyPayload struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	return func(ctx *gin.Context) {
		var bodyData bodyPayload

		err := ctx.BindJSON(&bodyData)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Missing required body fields"})
			ctx.Abort()
			return
		}

		user, err := controllers.RetreiveUserByEmail(bodyData.Email)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid credentials"})
			ctx.Abort()
			return
		}

		err = utils.CheckPassword(bodyData.Password, user.Password)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid credentials"})
			ctx.Abort()
			return
		}

		access_token, err := controllers.GenerateToken(user.Email, 24*time.Hour)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Somthing went wrong"})
			ctx.Abort()
			return
		}
		refresh_token, err := controllers.GenerateToken(user.Email, 30*24*time.Hour)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Somthing went wrong"})
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
	type bodyPayload struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	return func(ctx *gin.Context) {
		var bodyData bodyPayload

		err := ctx.BindJSON(&bodyData)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Missing required body fields"})
			ctx.Abort()
			return
		}

		email, err := controllers.ValidateToken(bodyData.RefreshToken)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid auth token"})
			ctx.Abort()
			return
		}

		user, err := controllers.RetreiveUserByEmail(email)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user"})
			ctx.Abort()
			return
		}

		accessToken, err := controllers.GenerateToken(user.Email, 24*time.Hour)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
			ctx.Abort()
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message":       "Refreshed successfully",
			"access_token":  accessToken,
			"refresh_token": bodyData.RefreshToken,
		})
	}
}
