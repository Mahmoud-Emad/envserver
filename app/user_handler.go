package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"

	internal "github.com/Mahmoud-Emad/envserver/internal"
	models "github.com/Mahmoud-Emad/envserver/models"
)

// CreateUserInputs struct for data needed when user create an account
type CreateUserInputs struct {
	Name         string `json:"name" binding:"required" validate:"min=3,max=20"`
	Email        string `json:"email" binding:"required" validate:"mail"`
	Password     string `json:"password" binding:"required" validate:"password"`
	ProjectOwner bool   `json:"is_owner"`
}

// ValidateUser checks for the presence of required fields in the user struct.
func validateUserFields(user *CreateUserInputs) error {
	t := reflect.TypeOf(*user)
	v := reflect.ValueOf(*user)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Name != "ProjectOwner" {
			value := v.Field(i).Interface()
			// Check if the field value is empty or zero
			if reflect.DeepEqual(value, reflect.Zero(field.Type).Interface()) {
				return errors.New(fmt.Sprintf("%s field is required", field.Name))
			}
		}
	}

	return nil
}

// Get all users from the database handler.
func (a *App) getUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, _ := a.DB.GetUsers()
	sendJSONResponse(w, http.StatusOK, "Users found", users, nil)
	return
}

func (a *App) createUserHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request data
	var fields CreateUserInputs
	err := json.NewDecoder(r.Body).Decode(&fields)

	if err != nil {
		sendJSONResponse(w, http.StatusBadRequest, "Invalid request payload", nil, err)
		return
	}

	// Validate user data
	err = validateUserFields(&fields)
	if err != nil {
		sendJSONResponse(
			w, http.StatusBadRequest,
			"Please ensure that all mandatory fields have been filled out.",
			nil,
			err,
		)
		return
	}

	// We save the user password in the hash, the value should be the actual value and the key should be the hashed password value.
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

	found, err := a.DB.GetUserByEmail(fields.Email)
	if found.Email == fields.Email {
		sendJSONResponse(
			w, http.StatusBadRequest,
			"Failed to create user object.",
			nil,
			internal.UserEmailNotUniqueError,
		)
		return
	}

	user := models.User{
		Name:           fields.Name,
		Email:          fields.Email,
		HashedPassword: hashedPassword,
		Projects:       []*models.Project{},
		IsOwner:        fields.ProjectOwner,
	}

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
