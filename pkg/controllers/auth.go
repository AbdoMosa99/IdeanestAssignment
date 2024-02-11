package controllers

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = "secret"

// Generates JWT token with email string as the subject
// (which is the unique value for the token holder)
// and adds duration to current time as the expiration date
// then sign the token with our very secret key then return its string representation.
func GenerateToken(email string, duration time.Duration) (string, error) {
	// prepare the token with claims
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss": "AuthService",                           // Issuer
			"sub": email,                                   // Subject
			"exp": time.Now().Local().Add(duration).Unix(), // Expiration Date
		})

	// sign the token and return
	signedToken, err := token.SignedString([]byte(secretKey))
	return signedToken, err
}

// Validate the given signed token for that
//  1. we have signed it (so it's sealed)
//  2. have a subject claim holding the email
//  3. not expired yet
//
// then we return the underlying email
func ValidateToken(signedToken string) (string, error) {
	// unsign the token and make sure its ours
	token, err := jwt.Parse(
		signedToken,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})
	if err != nil {
		return "", err
	}

	// extract the email subject from it
	email, err := token.Claims.GetSubject()
	if err != nil {
		return "", err
	}

	// check that it's not expired yet
	expiresAt, err := token.Claims.GetExpirationTime()
	if expiresAt.Unix() < time.Now().Local().Unix() {
		err = errors.New("token is expired")
	}

	return email, err
}
