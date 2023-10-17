package app

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	models "github.com/Mahmoud-Emad/envserver/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/rs/zerolog/log"
)

type Response struct {
	Message string      `json:"message"`
	Status  int         `json:"status"`
	Data    interface{} `json:"data"`
}

type Handler func(w http.ResponseWriter, r *http.Request)

func sendJSONResponse(w http.ResponseWriter, status int, message string, data interface{}, err error) {

	response := Response{
		Status:  status,
		Message: message,
		Data:    data,
	}

	if err != nil {
		response.Message = message + ": " + err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

// wrapRequest wraps an HTTP handler function with additional debugging information and optional authentication check.
// If protected is true, the incoming request is expected to include a JWT token in the "Authorization" header.
func (a *App) wrapRequest(h http.HandlerFunc, protected bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if protected {
			// Check if the request includes a JWT token in the "Authorization" header.
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				sendJSONResponse(w, http.StatusUnauthorized, "Unauthorized: JWT token missing", nil, nil)
				return
			}

			// Validate and decode the JWT token.
			user, err := a.VerifyAndDecodeJwtToken(authHeader, a.Config.Server.JWTSecretKey)
			if err != nil {
				sendJSONResponse(w, http.StatusUnauthorized, "Unauthorized: Invalid JWT token", nil, err)
				return
			}

			// Add the user object inside the request.
			ctx := context.WithValue(r.Context(), UserContextKey, user)
			r = r.WithContext(ctx)

			// Print out the incoming request URL and its method for debugging.
			log.Debug().Msgf("User: %d | Request: %s %s", user.ID, r.Method, r.URL.Path)
		} else {
			// Print out the incoming request URL and its method for debugging.
			log.Debug().Msgf("Request: %s %s", r.Method, r.URL.Path)
		}
		// Call the original handler.
		h(w, r)
	}
}

func (a *App) authenticateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the JWT token from the Authorization header
		tokenString := r.Header.Get("Authorization")
		if strings.TrimSpace(tokenString) == "" {
			log.Warn().Msgf("Request|unauthorized: %s %s", r.Method, r.URL.Path)
			sendJSONResponse(w, http.StatusUnauthorized, "Authorization token required", nil, nil)
			return
		}

		// Parse the JWT token and validate its signature
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Replace "your-secret-key" with the same secret key used for signing the token
			return []byte(a.Config.Server.JWTSecretKey), nil
		})

		if err != nil || !token.Valid {
			log.Warn().Msgf("Request|unauthorized: %s %s", r.Method, r.URL.Path)
			sendJSONResponse(w, http.StatusUnauthorized, "Invalid or expired token", nil, err)
			return
		}

		// Token is valid, continue with the next handler
		next.ServeHTTP(w, r)
	})
}

type contextKey string

const UserContextKey contextKey = "user"

// Get the requested user data.
func (a *App) GetRequestedUser(r *http.Request) (models.User, error) {
	user, ok := r.Context().Value(UserContextKey).(models.User)
	if !ok {
		authHeader := r.Header.Get("Authorization")
		user, err := a.VerifyAndDecodeJwtToken(authHeader, a.Config.Server.JWTSecretKey)
		if err != nil {
			return user, errors.New("cannot decode jwt")
		}

		ctx := context.WithValue(r.Context(), UserContextKey, user)
		r.WithContext(ctx)
	}
	return user, nil
}
