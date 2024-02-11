package controllers

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = "secret"

func GenerateToken(email string, duration time.Duration) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss": "AuthService",
			"sub": email,
			"exp": time.Now().Local().Add(duration).Unix(),
		})
	signedToken, err := token.SignedString([]byte(secretKey))
	return signedToken, err
}

func ValidateToken(signedToken string) (string, error) {
	token, err := jwt.Parse(
		signedToken,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})
	if err != nil {
		return "", err
	}

	email, err := token.Claims.GetSubject()
	if err != nil {
		return "", err
	}

	expiresAt, err := token.Claims.GetExpirationTime()
	if expiresAt.Unix() < time.Now().Local().Unix() {
		err = errors.New("token is expired")
	}

	return email, err
}
