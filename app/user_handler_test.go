package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Mahmoud-Emad/envserver/internal"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var configContent = `
[server]
host = "localhost"
port = 8080
jwt_secret_key = "xyz"
shutdown_timeout = 10

[database]
host = "localhost"
port = 5432
name = "postgres"
user = "postgres"
password = "postgres"
`

var userToken = ""

// Create a temporary file
func createConfTempFile(t *testing.T) *os.File {
	tempFile, err := os.CreateTemp("", "config.toml")
	assert.NoError(t, err)

	err = os.WriteFile(tempFile.Name(), []byte(configContent), 0644)
	assert.NoError(t, err)
	return tempFile
}

func getUserToken(t *testing.T, responseRecorder *httptest.ResponseRecorder) string {
	var responseBody map[string]interface{}
	err := json.NewDecoder(responseRecorder.Body).Decode(&responseBody)
	assert.NoError(t, err)

	// Check if the "data" field exists and is a map
	data, found := responseBody["data"].(map[string]interface{})
	assert.True(t, found, "Data field not found in the response body")

	// Extract the token from the "data" field
	token, found := data["token"].(string)
	assert.True(t, found, "Token not found in the data field")
	assert.NotEmpty(t, token)

	return token
}

func TestDeleteUserByIDHandler(t *testing.T) {
	tempFile := createConfTempFile(t)
	registerEndpoint := "/api/v1/auth/signin"
	loginEndpoint := "/api/v1/auth/signin"

	// Close the temporary file after creating the App instance
	defer func() {
		tempFile.Close()
		os.Remove(tempFile.Name())
	}()

	t.Run("Success registration", func(t *testing.T) {
		user := internal.SignUpInputs{
			Name:         "omda",
			Email:        "omda@test.delete",
			Password:     "password123",
			ProjectOwner: false,
		}

		jsonPayload, err := json.Marshal(user)
		assert.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, registerEndpoint, strings.NewReader(string(jsonPayload)))
		request.Header.Set("Content-Type", "application/json")

		app, err := NewApp(tempFile.Name())
		assert.NoError(t, err)

		responseRecorder := httptest.NewRecorder()
		app.signupHandler(responseRecorder, request)
		assert.Equal(t, responseRecorder.Result().StatusCode, http.StatusCreated)
	})

	t.Run("Success loggedin", func(t *testing.T) {
		user := internal.SignUpInputs{
			Email:    "omda@test.delete",
			Password: "password123",
		}

		jsonPayload, err := json.Marshal(user)
		assert.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, loginEndpoint, strings.NewReader(string(jsonPayload)))
		request.Header.Set("Content-Type", "application/json")

		app, err := NewApp(tempFile.Name())
		assert.NoError(t, err)

		responseRecorder := httptest.NewRecorder()
		app.signinHandler(responseRecorder, request)
		assert.Equal(t, responseRecorder.Result().StatusCode, http.StatusOK)

		userToken = getUserToken(t, responseRecorder)
		assert.NotEmpty(t, userToken)
	})

	t.Run("Success delete user", func(t *testing.T) {
		app, err := NewApp(tempFile.Name())
		assert.NoError(t, err)

		user, err := app.VerifyAndDecodeJwtToken(userToken, app.Config.Server.JWTSecretKey)
		assert.NoError(t, err)
		assert.NotEqual(t, user.ID, 0)

		// Construct the URL with the correct user ID
		url := fmt.Sprintf("/api/v1/users/%d", user.ID)

		// Create the DELETE request with the user token
		request := httptest.NewRequest(http.MethodDelete, url, nil)
		request = mux.SetURLVars(request, map[string]string{"id": fmt.Sprint(user.ID)})

		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", userToken)

		responseRecorder := httptest.NewRecorder()
		app.deleteUserByIDHandler(responseRecorder, request)

		// Expecting a successful deletion with status code 204
		assert.Equal(t, responseRecorder.Result().StatusCode, http.StatusNoContent)
	})

}
