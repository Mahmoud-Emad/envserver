package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	internal "github.com/Mahmoud-Emad/envserver/internal"
	"github.com/stretchr/testify/assert"
)

func TestSignupHandler(t *testing.T) {
	tempFile := createConfTempFile(t)
	registerEndpoint := "/api/v1/auth/signup"

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

		// Delete the created user.
		app.DB.DeleteUserByEmail(user.Email)
		_, err = app.DB.GetUserByEmail(user.Email)
		assert.Error(t, err)
	})

	t.Run("Fail registration <Name is missing>", func(t *testing.T) {
		user := internal.SignUpInputs{
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
		assert.Equal(t, responseRecorder.Result().StatusCode, http.StatusBadRequest)
	})
}

func TestSigninHandler(t *testing.T) {
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
	})

	t.Run("login fail", func(t *testing.T) {
		user := internal.SignUpInputs{
			Email:    "omda@gmail.com",
			Password: "1245sa",
		}

		jsonPayload, err := json.Marshal(user)
		assert.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, loginEndpoint, strings.NewReader(string(jsonPayload)))
		request.Header.Set("Content-Type", "application/json")

		app, err := NewApp(tempFile.Name())
		assert.NoError(t, err)

		responseRecorder := httptest.NewRecorder()
		app.signinHandler(responseRecorder, request)
		assert.Equal(t, responseRecorder.Result().StatusCode, http.StatusUnauthorized)
	})

	t.Run("delete the created user", func(t *testing.T) {
		user := internal.SignUpInputs{
			Email: "omda@gmail.com",
		}

		app, err := NewApp(tempFile.Name())
		assert.NoError(t, err)

		app.DB.DeleteUserByEmail(user.Email)
		_, err = app.DB.GetUserByEmail(user.Email)
		assert.Error(t, err)
	})
}
