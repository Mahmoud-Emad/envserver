package app

import (
	"encoding/json"
	"net/http"
	"strconv"

	internal "github.com/Mahmoud-Emad/envserver/internal"
	"github.com/Mahmoud-Emad/envserver/models"
	"github.com/gorilla/mux"
)

var envFields internal.EnvironmentKeyInputs

// getProjectEnvHandler retrieves a list of project env vars from the database and sends the response as JSON.
func (a *App) getProjectEnvHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if len(vars) == 0 {
		sendJSONResponse(
			w,
			http.StatusBadRequest,
			"Cannot find the project by the provided ID.",
			nil,
			internal.ProjectIdNotProvidedError,
		)
		return
	}

	projectIDStr := vars["id"]
	convertedProjectId, err := strconv.ParseInt(projectIDStr, 10, 64)

	if err != nil {
		sendJSONResponse(
			w,
			http.StatusBadRequest,
			"Cannot convert project id to number.",
			nil,
			err,
		)
	}

	pId := int(convertedProjectId)
	env, err := a.DB.GetEnvKeysAndValuesById(pId)

	if err != nil {
		sendJSONResponse(
			w,
			http.StatusNotFound,
			"Failed to retrieve project environment",
			nil,
			err,
		)
		return
	}
	sendJSONResponse(w, http.StatusOK, "Project environment found successfully", env, nil)
}

// Create new env key/value inside a project
func (a *App) createProjectEnvHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if len(vars) == 0 {
		sendJSONResponse(
			w,
			http.StatusBadRequest,
			"The project id should be provided.",
			nil,
			internal.ProjectIdNotProvidedError,
		)
		return
	}

	projectIDStr := vars["id"]
	convertedProjectId, err := strconv.ParseInt(projectIDStr, 10, 64)

	if err != nil {
		sendJSONResponse(
			w,
			http.StatusBadRequest,
			"Cannot convert project id to number.",
			nil,
			err,
		)
		return
	}

	_, err = a.DB.GetProjectByID(int(convertedProjectId))
	if err != nil {
		sendJSONResponse(
			w,
			http.StatusBadRequest,
			"The project does not exist",
			nil,
			err,
		)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&envFields)
	if err != nil {
		sendJSONResponse(
			w,
			http.StatusBadRequest,
			"Invalid request payload",
			nil,
			err,
		)
		return
	}

	// Validate user data
	err = envFields.Validate()
	if err != nil {
		sendJSONResponse(
			w, http.StatusBadRequest,
			"Please ensure that all mandatory fields have been filled out.",
			nil,
			err,
		)
		return
	}

	hashedValue, err := internal.HashPassword(envFields.Value)
	if err != nil {
		sendJSONResponse(
			w,
			http.StatusBadRequest,
			"Error hashing value:",
			nil,
			err,
		)
		return
	}

	// Create new project object.
	env := models.EnvironmentKey{
		Key:       envFields.Key,
		Value:     hashedValue,
		ProjectID: int(convertedProjectId),
	}

	err = a.DB.CreateEnvKey(&env)
	if err != nil {
		sendJSONResponse(
			w, http.StatusBadRequest,
			"Failed to create project environment object",
			nil,
			err,
		)
		return
	}

	envFields = internal.EnvironmentKeyInputs{}
	sendJSONResponse(w, http.StatusCreated, "Project environment created successfully", env, nil)
}
