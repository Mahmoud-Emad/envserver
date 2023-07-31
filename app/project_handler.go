package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	internal "github.com/Mahmoud-Emad/envserver/internal"
	"github.com/Mahmoud-Emad/envserver/models"
	"github.com/gorilla/mux"
)

func (a *App) getProjectsHandler(w http.ResponseWriter, r *http.Request) {
	projects, err := a.DB.GetProjects()
	if err != nil {
		sendJSONResponse(w, http.StatusBadRequest, "Failed to retrieve projects", nil, err)
		return
	}
	sendJSONResponse(w, http.StatusOK, "Projects found successfully", projects, nil)
}

func (a *App) deleteProjectByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if len(vars) == 0 {
		sendJSONResponse(w, http.StatusBadRequest, "Cannot get the project id.", nil, errors.New("Project id should be provided."))
	}

	projectIDStr := vars["id"]
	// Convert the user ID to uint
	u64, err := strconv.ParseUint(projectIDStr, 10, 32)
	uID := uint(u64)

	project, err := a.DB.GetProjectByID(uID)
	if err != nil {
		sendJSONResponse(w, http.StatusNotFound, fmt.Sprintf("Failed to retrieve project with id %s.", projectIDStr), nil, err)
		return
	}

	a.DB.DeleteProjectByID(project.ID)
	sendJSONResponse(w, http.StatusNoContent, "Project deleted successfully", nil, nil)
}

func (a *App) getProjectByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if len(vars) == 0 {
		sendJSONResponse(w, http.StatusBadRequest, "Cannot get the project id.", nil, errors.New("Project id should be provided."))
	}

	projectIDStr := vars["id"]
	// Convert the user ID to uint
	u64, err := strconv.ParseUint(projectIDStr, 10, 32)
	uID := uint(u64)

	project, err := a.DB.GetProjectByID(uID)
	if err != nil {
		sendJSONResponse(w, http.StatusNotFound, fmt.Sprintf("Failed to retrieve project with id %s.", projectIDStr), nil, err)
		return
	}
	sendJSONResponse(w, http.StatusOK, "Project found successfully", project, nil)
}

func (a *App) createProjectHandler(w http.ResponseWriter, r *http.Request) {
	user, err := internal.GetRequestedUser(r)
	if err != nil {
		sendJSONResponse(w, http.StatusBadRequest, "Requested user not found.", nil, err)
		return
	}

	var fields internal.ProjectInputs
	err = json.NewDecoder(r.Body).Decode(&fields)

	if err != nil {
		sendJSONResponse(w, http.StatusBadRequest, "Invalid request payload", nil, err)
		return
	}

	// Validate user data
	err = internal.ValidateProjectFields(&fields)
	if err != nil {
		sendJSONResponse(
			w, http.StatusBadRequest,
			"Please ensure that all mandatory fields have been filled out.",
			nil,
			err,
		)
		return
	}

	// Create new project object.
	project := models.Project{
		Name:  fields.Name,
		Owner: user.ID,
		Team:  []*models.User{},
		Keys:  []*models.EnvironmentKey{},
	}

	// save the project into the database
	err = a.DB.CreateProject(&project)
	if err != nil {
		sendJSONResponse(
			w, http.StatusBadRequest,
			"Failed to create project object.",
			nil,
			err,
		)
		return
	}
	// Return success response
	sendJSONResponse(w, http.StatusCreated, "Project created successfully", project, nil)
	// TODO, Make the project unique on the user and the env type -> [test, dev..etc]
}
