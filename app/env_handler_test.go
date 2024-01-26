package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	internal "github.com/Mahmoud-Emad/envserver/internal"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestProjectEnv(t *testing.T) {
	tempFile := createConfTempFile(t)
	projectsEndpoint := "/api/v1/projects"
	registerEndpoint := "/api/v1/auth/signup"
	loginEndpoint := "/api/v1/auth/signin"
	projectID := ""

	// Close the temporary file after creating the App instance
	defer func() {
		tempFile.Close()
		os.Remove(tempFile.Name())
	}()

	t.Run("Success registration", func(t *testing.T) {
		user := internal.SignUpInputs{
			FirstName:    "omda",
			LastName:     "Man",
			Email:        "omda@env.com",
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
			Email:    "omda@env.com",
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
	})

	t.Run("Test success create project", func(t *testing.T) {
		app, err := NewApp(tempFile.Name())
		assert.NoError(t, err)

		user, err := app.VerifyAndDecodeJwtToken(userToken, app.Config.Server.JWTSecretKey)
		assert.NoError(t, err)
		assert.NotEqual(t, user.ID, 0)

		jsonPayload, err := json.Marshal(projectData)
		assert.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, projectsEndpoint, strings.NewReader(string(jsonPayload)))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", userToken)

		responseRecorder := httptest.NewRecorder()
		app.createProjectHandler(responseRecorder, request)
		assert.Equal(t, responseRecorder.Result().StatusCode, http.StatusCreated)

		projectID = getProjectID(t, responseRecorder)
	})

	t.Run("Test success get project env", func(t *testing.T) {
		app, err := NewApp(tempFile.Name())
		assert.NoError(t, err)

		user, err := app.VerifyAndDecodeJwtToken(userToken, app.Config.Server.JWTSecretKey)
		assert.NoError(t, err)
		assert.NotEqual(t, user.ID, 0)

		projectsEndpoint = fmt.Sprintf("/api/v1/projects/%s/env", projectID)

		request := httptest.NewRequest(http.MethodGet, projectsEndpoint, nil)
		request = mux.SetURLVars(request, map[string]string{"id": fmt.Sprint(projectID)})

		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", userToken)

		responseRecorder := httptest.NewRecorder()
		app.getProjectEnvHandler(responseRecorder, request)

		var responseBody map[string]interface{}
		err = json.NewDecoder(responseRecorder.Body).Decode(&responseBody)
		assert.NoError(t, err)
		assert.Equal(t, responseRecorder.Result().StatusCode, http.StatusOK)
	})

	t.Run("Test success delete project", func(t *testing.T) {
		app, err := NewApp(tempFile.Name())
		assert.NoError(t, err)

		user, err := app.VerifyAndDecodeJwtToken(userToken, app.Config.Server.JWTSecretKey)
		assert.NoError(t, err)
		assert.NotEqual(t, user.ID, 0)

		projectsEndpoint = "/api/v1/projects"

		// Construct the URL with the correct user ID
		url := fmt.Sprintf("%s/%s", projectsEndpoint, projectID)

		request := httptest.NewRequest(http.MethodDelete, url, nil)
		request = mux.SetURLVars(request, map[string]string{"id": fmt.Sprint(projectID)})

		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", userToken)

		responseRecorder := httptest.NewRecorder()
		app.deleteProjectByIDHandler(responseRecorder, request)
		assert.Equal(t, responseRecorder.Result().StatusCode, http.StatusNoContent)
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

// Get the project id from the response
func getProjectID(t *testing.T, responseRecorder *httptest.ResponseRecorder) string {
	var responseBody map[string]interface{}
	err := json.NewDecoder(responseRecorder.Body).Decode(&responseBody)
	assert.NoError(t, err)

	// Check if the "data" field exists and is a map
	data, found := responseBody["data"].(map[string]interface{})
	assert.True(t, found, "Data field not found in the response body")

	projectID, found := data["ID"].(float64)
	assert.True(t, found, "ID field not found in the response body")

	sID := fmt.Sprintf("%d", int(projectID)) // convert to string
	return sID
}
