package app

import (
	"errors"

	models "github.com/Mahmoud-Emad/envserver/models"
	"github.com/dgrijalva/jwt-go"
)

func GenerateJwtToken(payload map[string]interface{}, JWTSecretKey string) string {
	// Generate a JWT token with user data as the payload
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(payload))

	tokenString, err := token.SignedString([]byte(JWTSecretKey))

	if err != nil {
		panic(err)
	}

	return tokenString
}

func VerifyAndDecodeJwtToken(tokenString, JWTSecretKey string) (models.User, error) {
	// Parse the token and extract the payload.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWTSecretKey), nil
	})

	if err != nil {
		return models.User{}, err
	}

	if !token.Valid {
		return models.User{}, errors.New("invalid token")
	}

	// Extract the payload data
	payload, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return models.User{}, errors.New("invalid token claims")
	}

	user := models.User{
		ID:    payload["id"].(uint),
		Email: payload["email"].(string),
	}
	return user, nil
}
