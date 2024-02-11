package handlers

import (
	"ideanest/pkg/controllers"
	"ideanest/pkg/database/mongodb/models"
	"ideanest/pkg/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Register a new user
func SignUp() gin.HandlerFunc {
	// body format
	type bodyInterface struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	return func(ctx *gin.Context) {
		// check that body matches our format
		var bodyData bodyInterface
		err := ctx.BindJSON(&bodyData)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Missing required body fields"})
			return
		}

		// create the user into the database
		user := models.UserModel{
			Name:     bodyData.Name,
			Email:    bodyData.Email,
			Password: utils.HashPassword(bodyData.Password),
		}
		err = controllers.InsertUser(user)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
			return
		}

		// return success
		ctx.JSON(http.StatusOK, gin.H{"message": "Registered Successfully"})
	}
}

// Login a user
func SignIn() gin.HandlerFunc {
	// body format
	type bodyInterface struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	return func(ctx *gin.Context) {
		// check that body matches our format
		var bodyData bodyInterface
		err := ctx.BindJSON(&bodyData)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Missing required body fields"})
			return
		}

		// get that user if exists
		user, err := controllers.RetreiveUserByEmail(bodyData.Email)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid credentials"})
			return
		}

		// check its password
		err = utils.CheckPassword(bodyData.Password, user.Password)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid credentials"})
			return
		}

		// generate an access token (valid for 24 hrs)
		access_token, err := controllers.GenerateToken(user.Email, 24*time.Hour)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Somthing went wrong"})
			return
		}

		// generate a refresh token (valid for a month)
		refresh_token, err := controllers.GenerateToken(user.Email, 30*24*time.Hour)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Somthing went wrong"})
			return
		}

		// return success
		ctx.JSON(http.StatusOK, gin.H{
			"message":       "Logged in successfully",
			"access_token":  access_token,
			"refresh_token": refresh_token,
		})
	}
}

// refresh user token using his refresh token
func RefreshToken() gin.HandlerFunc {
	// body format
	type bodyInterface struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	return func(ctx *gin.Context) {
		// check that body matches our format
		var bodyData bodyInterface
		err := ctx.BindJSON(&bodyData)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Missing required body fields"})
			return
		}

		// make sure the refresh token provided is valid
		email, err := controllers.ValidateToken(bodyData.RefreshToken)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid auth token"})
			return
		}

		// get that user
		user, err := controllers.RetreiveUserByEmail(email)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user"})
			return
		}

		// generate a new token for 24 hrs
		accessToken, err := controllers.GenerateToken(user.Email, 24*time.Hour)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
			return
		}

		// return success
		ctx.JSON(http.StatusOK, gin.H{
			"message":       "Refreshed successfully",
			"access_token":  accessToken,
			"refresh_token": bodyData.RefreshToken,
		})
	}
}
