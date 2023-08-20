package app

import (
	"encoding/json"
	"net/http"

	internal "github.com/Mahmoud-Emad/envserver/internal"
	models "github.com/Mahmoud-Emad/envserver/models"
)

func (a *App) signinHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request data
	var fields internal.SigninInputs

	if err := json.NewDecoder(r.Body).Decode(&fields); err != nil {
		sendJSONResponse(w, http.StatusBadRequest, "Invalid request payload", nil, err)
		return
	}

	// Find the user by email
	user, err := a.DB.GetUserByEmail(fields.Email)
	if err != nil {
		sendJSONResponse(w, http.StatusUnauthorized, "Cannot get user object with this email.", nil, err)
		return
	}

	// Check if the provided password is correct
	if !internal.CheckPasswordHash(fields.Password, user.HashedPassword) {
		sendJSONResponse(w, http.StatusUnauthorized, "Invalid email or password", nil, nil)
		return
	}

	// Generate a JWT token with user data as the payload
	payload := map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
	}

	token, err := GenerateJwtToken(payload, a.Config.Server.JWTSecretKey)

	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, "Failed to generate JWT token", nil, err)
		return
	}

	// Return success response with JWT token
	sendJSONResponse(w, http.StatusOK, "User authenticated successfully", map[string]string{"token": token}, nil)
}

// signupHandler handles the HTTP request for creating a user.
// It expects a JSON payload containing user data in the request body.
// If the request is valid and the user is created successfully, it returns a JSON response with status 201 (Created).
// If the request is invalid or encounters an error, it returns an appropriate error response.
func (a *App) signupHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request data
	var fields internal.SignUpInputs
	err := json.NewDecoder(r.Body).Decode(&fields)

	if err != nil {
		sendJSONResponse(w, http.StatusBadRequest, "Invalid request payload", nil, err)
		return
	}

	// Validate user data
	err = internal.ValidateUserFields(&fields)
	if err != nil {
		sendJSONResponse(
			w, http.StatusBadRequest,
			"Please ensure that all mandatory fields have been filled out.",
			nil,
			err,
		)
		return
	}

	// Hash the user password
	mdPassHash := internal.HashMD5(fields.Password)
	hashedPassword, err := internal.EncryptAES([]byte(fields.Password), mdPassHash)

	if err != nil {
		sendJSONResponse(
			w, http.StatusBadRequest,
			"Failed to create user object.",
			nil,
			internal.InternalServerError,
		)
		return
	}

	// Check if the email is already taken
	found, _ := a.DB.GetUserByEmail(fields.Email)

	if found.Email == fields.Email {
		sendJSONResponse(
			w, http.StatusBadRequest,
			"Failed to create user object.",
			nil,
			internal.UserEmailNotUniqueError,
		)
		return
	}

	// Create the user object
	user := models.User{
		Name:           fields.Name,
		Email:          fields.Email,
		HashedPassword: hashedPassword,
		Projects:       []*models.Project{},
		IsOwner:        fields.ProjectOwner,
	}

	// Save the user in the database
	err = a.DB.CreateUser(&user)
	if err != nil {
		sendJSONResponse(
			w, http.StatusBadRequest,
			"Failed to create user object.",
			nil,
			err,
		)
		return
	}

	// Return success response
	sendJSONResponse(w, http.StatusCreated, "User registered successfully", user, nil)
}
