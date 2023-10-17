package app

import (
	"errors"
	"fmt"

	models "github.com/Mahmoud-Emad/envserver/models"
	"github.com/dgrijalva/jwt-go"
)

func (a *App) GenerateJwtToken(payload map[string]interface{}, JWTSecretKey string) (string, error) {
	// Generate a JWT token with user data as the payload
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(payload))

	tokenString, err := token.SignedString([]byte(JWTSecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (a *App) VerifyAndDecodeJwtToken(tokenString, JWTSecretKey string) (models.User, error) {
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

	idFloat, ok := payload["id"].(float64)
	if !ok {
		return models.User{}, fmt.Errorf("id %v is an invalid id field in token", payload["id"])
	}

	id := int(idFloat)

	user, err := a.DB.GetUserByID(id)
	if err != nil {
		return models.User{}, fmt.Errorf("cannot find a user with id %d", id)
	}

	return user, nil
}
