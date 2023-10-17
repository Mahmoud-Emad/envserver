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
	models "github.com/Mahmoud-Emad/envserver/models"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var projectData = internal.ProjectInputs{
	Name: "testProject",
}

func TestProjectHandlers(t *testing.T) {
	tempFile := createConfTempFile(t)
	registerEndpoint := "/api/v1/auth/signin"
	loginEndpoint := "/api/v1/auth/signin"
	projectsEndpoint := "/api/v1/projects"
	fakeID := 500000000221144

	// Close the temporary file after creating the App instance
	defer func() {
		tempFile.Close()
		os.Remove(tempFile.Name())
	}()

	t.Run("Success registration", func(t *testing.T) {
		user := internal.SignUpInputs{
			Name:         "omda",
			Email:        "omda@gmail.com",
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
			Email:    "omda@gmail.com",
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
	})

	t.Run("Test success get project", func(t *testing.T) {
		app, err := NewApp(tempFile.Name())
		assert.NoError(t, err)

		user, err := app.VerifyAndDecodeJwtToken(userToken, app.Config.Server.JWTSecretKey)
		assert.NoError(t, err)
		assert.NotEqual(t, user.ID, 0)

		project, err := app.DB.GetProjectByName(projectData.Name)
		assert.NoError(t, err)

		// Construct the URL with the correct user ID
		url := fmt.Sprintf("/%s/%d", projectsEndpoint, project.ID)

		request := httptest.NewRequest(http.MethodGet, url, nil)
		request = mux.SetURLVars(request, map[string]string{"id": fmt.Sprint(project.ID)})

		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", userToken)

		responseRecorder := httptest.NewRecorder()
		app.getProjectByIDHandler(responseRecorder, request)
		assert.Equal(t, responseRecorder.Result().StatusCode, http.StatusOK)
	})

	t.Run("Test fail to get project", func(t *testing.T) {
		app, err := NewApp(tempFile.Name())
		assert.NoError(t, err)

		user, err := app.VerifyAndDecodeJwtToken(userToken, app.Config.Server.JWTSecretKey)
		assert.NoError(t, err)
		assert.NotEqual(t, user.ID, 0)

		// Construct the URL with the correct user ID
		url := fmt.Sprintf("/%s/%d", projectsEndpoint, fakeID)

		request := httptest.NewRequest(http.MethodGet, url, nil)
		request = mux.SetURLVars(request, map[string]string{"id": fmt.Sprint(fakeID)})

		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", userToken)

		responseRecorder := httptest.NewRecorder()
		app.getProjectByIDHandler(responseRecorder, request)
		assert.Equal(t, responseRecorder.Result().StatusCode, http.StatusNotFound)
	})

	t.Run("Test success update project", func(t *testing.T) {
		app, err := NewApp(tempFile.Name())
		assert.NoError(t, err)

		user, err := app.VerifyAndDecodeJwtToken(userToken, app.Config.Server.JWTSecretKey)
		assert.NoError(t, err)
		assert.NotEqual(t, user.ID, 0)

		p := models.Project{
			Name: "createProjectForTest",
		}

		err = app.DB.CreateProject(&p)
		assert.NoError(t, err)

		project, err := app.DB.GetProjectByName(p.Name)
		assert.NoError(t, err)

		p.Name = "testUpdatedProject"
		url := fmt.Sprintf("/%s/%d", projectsEndpoint, project.ID)

		jsonPayload, err := json.Marshal(p)
		assert.NoError(t, err)

		request := httptest.NewRequest(http.MethodPut, url, strings.NewReader(string(jsonPayload)))
		request = mux.SetURLVars(request, map[string]string{"id": fmt.Sprint(project.ID)})

		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", userToken)

		responseRecorder := httptest.NewRecorder()
		app.updateProjectHandler(responseRecorder, request)
		assert.Equal(t, responseRecorder.Result().StatusCode, http.StatusOK)

		var responseBody map[string]interface{}
		err = json.NewDecoder(responseRecorder.Body).Decode(&responseBody)
		assert.NoError(t, err)

		// Check if the "data" field exists and is a map
		data, found := responseBody["data"].(map[string]interface{})
		assert.True(t, found, "Data field not found in the response body")
		assert.Equal(t, data["name"], p.Name)

		// Delete created project.
		err = app.DB.DeleteProjectByName(p.Name)
		assert.NoError(t, err)
	})

	t.Run("Test fail to create project", func(t *testing.T) {
		// Empty name error.
		project := internal.ProjectInputs{}

		jsonPayload, err := json.Marshal(project)
		assert.NoError(t, err)
		request := httptest.NewRequest(http.MethodPost, projectsEndpoint, strings.NewReader(string(jsonPayload)))
		request.Header.Set("Content-Type", "application/json")

		app, err := NewApp(tempFile.Name())
		assert.NoError(t, err)

		responseRecorder := httptest.NewRecorder()
		app.createProjectHandler(responseRecorder, request)
		assert.Equal(t, responseRecorder.Result().StatusCode, http.StatusBadRequest)
	})

	// Delete created objects.
	t.Run("Test success delete project", func(t *testing.T) {
		app, err := NewApp(tempFile.Name())
		assert.NoError(t, err)

		user, err := app.VerifyAndDecodeJwtToken(userToken, app.Config.Server.JWTSecretKey)
		assert.NoError(t, err)
		assert.NotEqual(t, user.ID, 0)

		project, err := app.DB.GetProjectByName(projectData.Name)
		assert.NoError(t, err)

		// Construct the URL with the correct user ID
		url := fmt.Sprintf("/%s/%d", projectsEndpoint, project.ID)

		request := httptest.NewRequest(http.MethodDelete, url, nil)
		request = mux.SetURLVars(request, map[string]string{"id": fmt.Sprint(project.ID)})

		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", userToken)

		responseRecorder := httptest.NewRecorder()
		app.deleteProjectByIDHandler(responseRecorder, request)
		assert.Equal(t, responseRecorder.Result().StatusCode, http.StatusNoContent)
	})

	t.Run("Test fail to delete project", func(t *testing.T) {
		app, err := NewApp(tempFile.Name())
		assert.NoError(t, err)

		user, err := app.VerifyAndDecodeJwtToken(userToken, app.Config.Server.JWTSecretKey)
		assert.NoError(t, err)
		assert.NotEqual(t, user.ID, 0)

		// Construct the URL with the correct user ID
		url := fmt.Sprintf("/%s/%d", projectsEndpoint, fakeID)

		request := httptest.NewRequest(http.MethodDelete, url, nil)
		request = mux.SetURLVars(request, map[string]string{"id": fmt.Sprint(fakeID)})

		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", userToken)

		responseRecorder := httptest.NewRecorder()
		app.deleteProjectByIDHandler(responseRecorder, request)
		assert.Equal(t, responseRecorder.Result().StatusCode, http.StatusNotFound)

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
