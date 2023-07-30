package app

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

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

// Create a temporary file
func createConfTempFile(t *testing.T) *os.File {
	tempFile, err := ioutil.TempFile("", "config.toml")
	assert.NoError(t, err)

	err = ioutil.WriteFile(tempFile.Name(), []byte(configContent), 0644)
	assert.NoError(t, err)
	return tempFile
}

func requestToDeleteUserByID(t *testing.T, userID uint64) *http.Response {
	// Create a temporary file
	tempFile := createConfTempFile(t)
	userEndpoint := "/api/v1/users"

	// Close the temporary file after creating the App instance
	defer func() {
		tempFile.Close()
		os.Remove(tempFile.Name())
	}()

	// Create a test HTTP request
	url := fmt.Sprintf("%s/%d", userEndpoint, userID)
	request := httptest.NewRequest(http.MethodDelete, url, nil)
	request.Header.Set("Content-Type", "application/json")

	// Create a test HTTP response recorder
	responseRecorder := httptest.NewRecorder()

	// Call the handler function with the test request and response recorder
	app, err := NewApp(tempFile.Name()) // Use the temporary file as the config file
	assert.NoError(t, err)
	app.deleteUserByIDHandler(responseRecorder, request)

	// Check the HTTP response
	return responseRecorder.Result()
}
