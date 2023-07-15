package app

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Mahmoud-Emad/envserver/internal"
	"github.com/stretchr/testify/assert"
)

var configContent = `
[server]
host = "localhost"
port = 8080

[database]
host = "localhost"
port = 5432
name = "postgres"
user = "postgres"
password = "postgres"
`

var userObj = CreateUserInputs{
	Name:         "Mahmoud",
	Email:        "Mahmoud@gmail.com",
	Password:     "password123",
	ProjectOwner: false,
}

// Create a temporary file
func createConfTempFile(t *testing.T) *os.File {
	tempFile, err := ioutil.TempFile("", "config.toml")
	assert.NoError(t, err)

	err = ioutil.WriteFile(tempFile.Name(), []byte(configContent), 0644)
	assert.NoError(t, err)
	return tempFile
}

func TestCreateUserHandler(t *testing.T) {
	t.Run("success create user", func(t *testing.T) {
		response := requestToCreateUser(t, userObj)
		defer response.Body.Close()
		assert.Equal(t, http.StatusCreated, response.StatusCode)
	})

	t.Run("failed to create user (email not unique)", func(t *testing.T) {
		response := requestToCreateUser(t, userObj)
		defer response.Body.Close()

		assert.Equal(t, http.StatusBadRequest, response.StatusCode)

		var responseBody Response
		err := json.NewDecoder(response.Body).Decode(&responseBody)

		assert.NoError(t, err)
		assert.Equal(t, internal.UserEmailNotUniqueError.Error(), responseBody.Error.Message)
	})
}

// A helper function to make a request on the create user endpoint and then it returns the response.
func requestToCreateUser(t *testing.T, userData CreateUserInputs) *http.Response {
	// Create a temporary file
	tempFile := createConfTempFile(t)
	userEndPoint := "/api/v1/users"

	// Close the temporary file after creating the App instance
	defer func() {
		tempFile.Close()
		os.Remove(tempFile.Name())
	}()
	// Encode the User struct into JSON
	payload, err := json.Marshal(userData)

	// Create a test HTTP request
	request := httptest.NewRequest(http.MethodPost, userEndPoint, strings.NewReader(string(payload)))
	request.Header.Set("Content-Type", "application/json")

	// Create a test HTTP response recorder
	responseRecorder := httptest.NewRecorder()

	// Call the handler function with the test request and response recorder
	app, err := NewApp(tempFile.Name()) // Use the temporary file as the config file
	assert.NoError(t, err)
	app.createUserHandler(responseRecorder, request)

	// Check the HTTP response
	return responseRecorder.Result()
}
