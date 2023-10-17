package app

import (
	"encoding/json"
	"fmt"
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
			"The request variables have no items",
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
			"Cannot convert project id to number",
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

// updateProjectEnvKeyValueHandler is an endpoint to update the key/value of an exist key in the database by providing the object ID.
func (a *App) updateProjectEnvKeyValueHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if len(vars) == 0 {
		sendJSONResponse(
			w,
			http.StatusBadRequest,
			"The request variables have no items",
			nil,
			internal.ProjectIdNotProvidedError,
		)
	}

	projectIDStr := vars["projectID"]
	convertedProjectId, err := strconv.ParseInt(projectIDStr, 10, 64)

	if err != nil {
		sendJSONResponse(
			w,
			http.StatusBadRequest,
			"Cannot convert project id to number",
			nil,
			err,
		)
		return
	}

	_, err = a.DB.GetProjectByID(int(convertedProjectId))
	if err != nil {
		sendJSONResponse(
			w,
			http.StatusNotFound,
			fmt.Sprintf("Failed to retrieve project with id %s", projectIDStr),
			nil,
			err,
		)
		return
	}

	envIDStr := vars["envID"]
	convertedEnvId, err := strconv.ParseInt(envIDStr, 10, 64)

	if err != nil {
		sendJSONResponse(
			w,
			http.StatusBadRequest,
			"Cannot convert env object id to number",
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

	existingEnv, err := a.DB.GetProjectEnvByID(int(convertedEnvId))
	if err != nil {
		sendJSONResponse(w, http.StatusNotFound, fmt.Sprintf("Failed to retrieve project environment with id %s", envIDStr), nil, err)
		return
	}

	existingEnv.Key = envFields.Key
	existingEnv.Value = hashedValue

	err = a.DB.UpdateProjectEnvironment(existingEnv)
	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, "Failed to update project environment", nil, err)
		return
	}
	sendJSONResponse(w, http.StatusOK, "Project environment updated successfully", existingEnv, nil)
}

// Create new env key/value inside a project
func (a *App) createProjectEnvHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if len(vars) == 0 {
		sendJSONResponse(
			w,
			http.StatusBadRequest,
			"The request variables have no items",
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
			"Cannot convert project id to number",
			nil,
			err,
		)
		return
	}

	project, err := a.DB.GetProjectByID(int(convertedProjectId))
	if err != nil {
		sendJSONResponse(
			w,
			http.StatusNotFound,
			fmt.Sprintf("Failed to retrieve project with id %s", projectIDStr),
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
			w,
			http.StatusBadRequest,
			"Please ensure that all mandatory fields have been filled out",
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

	// Create new project environment object.
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

	// Update the project's keys.
	project.Keys = append(project.Keys, &env)
	err = a.DB.UpdateProject(&project)
	if err != nil {
		sendJSONResponse(
			w, http.StatusBadRequest,
			"Failed to update project object",
			nil,
			err,
		)
		return
	}

	envFields = internal.EnvironmentKeyInputs{}
	sendJSONResponse(w, http.StatusCreated, "Project environment created successfully", env, nil)
}

// getProjectEnvKeyValueHandler is an endpoint to get the key/value of an exist key in the database by providing the object ID.
func (a *App) getProjectEnvKeyValueHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if len(vars) == 0 {
		sendJSONResponse(
			w,
			http.StatusBadRequest,
			"The request variables have no items",
			nil,
			internal.ProjectIdNotProvidedError,
		)
	}

	projectIDStr := vars["projectID"]
	convertedProjectId, err := strconv.ParseInt(projectIDStr, 10, 64)

	if err != nil {
		sendJSONResponse(
			w,
			http.StatusBadRequest,
			"Cannot convert project id to number",
			nil,
			err,
		)
		return
	}

	_, err = a.DB.GetProjectByID(int(convertedProjectId))
	if err != nil {
		sendJSONResponse(
			w,
			http.StatusNotFound,
			fmt.Sprintf("Failed to retrieve project with id %s", projectIDStr),
			nil,
			err,
		)
		return
	}

	envIDStr := vars["envID"]
	convertedEnvId, err := strconv.ParseInt(envIDStr, 10, 64)

	if err != nil {
		sendJSONResponse(
			w,
			http.StatusBadRequest,
			"Cannot convert env object id to number",
			nil,
			err,
		)
		return
	}

	existingEnv, err := a.DB.GetProjectEnvByID(int(convertedEnvId))
	if err != nil {
		sendJSONResponse(w, http.StatusNotFound, fmt.Sprintf("Failed to retrieve project environment with id %s", envIDStr), nil, err)
		return
	}

	sendJSONResponse(w, http.StatusOK, "Project environment updated successfully", existingEnv, nil)
}

// deleteProjectEnvKeyValueHandler is an endpoint to delete the env object by providing the object ID.
func (a *App) deleteProjectEnvKeyValueHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if len(vars) == 0 {
		sendJSONResponse(
			w,
			http.StatusBadRequest,
			"The request variables have no items",
			nil,
			internal.ProjectIdNotProvidedError,
		)
	}

	projectIDStr := vars["projectID"]
	convertedProjectId, err := strconv.ParseInt(projectIDStr, 10, 64)

	if err != nil {
		sendJSONResponse(
			w,
			http.StatusBadRequest,
			"Cannot convert project id to number",
			nil,
			err,
		)
		return
	}

	_, err = a.DB.GetProjectByID(int(convertedProjectId))
	if err != nil {
		sendJSONResponse(
			w,
			http.StatusNotFound,
			fmt.Sprintf("Failed to retrieve project with id %s", projectIDStr),
			nil,
			err,
		)
		return
	}

	envIDStr := vars["envID"]
	convertedEnvId, err := strconv.ParseInt(envIDStr, 10, 64)

	if err != nil {
		sendJSONResponse(
			w,
			http.StatusBadRequest,
			"Cannot convert env object id to number",
			nil,
			err,
		)
		return
	}

	existingEnv, err := a.DB.GetProjectEnvByID(int(convertedEnvId))
	if err != nil {
		sendJSONResponse(w, http.StatusNotFound, fmt.Sprintf("Failed to retrieve project environment with id %s", envIDStr), nil, err)
		return
	}

	a.DB.DeleteProjectEnvByID(existingEnv.ID)
	sendJSONResponse(w, http.StatusNoContent, "Project environment deleted successfully", nil, nil)
}
