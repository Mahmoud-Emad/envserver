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

var projectFields internal.ProjectInputs

// getProjectsHandler retrieves a list of projects from the database and sends the response as JSON.
func (a *App) getProjectsHandler(w http.ResponseWriter, r *http.Request) {
	projects, err := a.DB.GetProjects()
	if err != nil {
		sendJSONResponse(w, http.StatusBadRequest, "Failed to retrieve projects", nil, err)
		return
	}
	sendJSONResponse(w, http.StatusOK, "Projects found successfully", projects, nil)
}

// deleteProjectByIDHandler deletes a project by its ID. The ID is retrieved from the URL path parameters.
func (a *App) deleteProjectByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if len(vars) == 0 {
		sendJSONResponse(w, http.StatusBadRequest, "Cannot get the project id.", nil, internal.ProjectIdNotProvidedError)
		return
	}

	projectIDStr := vars["id"]
	convertedProjectId, err := strconv.ParseInt(projectIDStr, 10, 64)

	if err != nil {
		sendJSONResponse(w, http.StatusBadRequest, "Cannot convert project id to number.", nil, err)
	}

	uId := int(convertedProjectId)

	project, err := a.DB.GetProjectByID(uId)
	if err != nil {
		sendJSONResponse(w, http.StatusNotFound, fmt.Sprintf("Failed to retrieve project with id %s.", projectIDStr), nil, err)
		return
	}

	a.DB.DeleteProjectByID(project.ID)
	sendJSONResponse(w, http.StatusNoContent, "Project deleted successfully", nil, nil)
}

// getProjectByIDHandler retrieves a project by its ID. The ID is retrieved from the URL path parameters.
func (a *App) getProjectByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if len(vars) == 0 {
		sendJSONResponse(w, http.StatusBadRequest, "Cannot get the project id.", nil, internal.ProjectIdNotProvidedError)
		return
	}

	projectIDStr := vars["id"]
	convertedProjectId, err := strconv.ParseInt(projectIDStr, 10, 64)

	if err != nil {
		sendJSONResponse(w, http.StatusBadRequest, "Cannot convert project id to number.", nil, err)
	}

	uId := int(convertedProjectId)

	project, err := a.DB.GetProjectByID(uId)
	if err != nil {
		sendJSONResponse(w, http.StatusNotFound, fmt.Sprintf("Failed to retrieve project with id %s.", projectIDStr), nil, err)
		return
	}
	sendJSONResponse(w, http.StatusOK, "Project found successfully", project, nil)
}

// createProjectHandler creates a new project based on the provided request payload.
func (a *App) createProjectHandler(w http.ResponseWriter, r *http.Request) {
	user, err := a.GetRequestedUser(r)
	if err != nil {
		sendJSONResponse(w, http.StatusBadRequest, "Requested user not found.", nil, err)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&projectFields)

	if err != nil {
		sendJSONResponse(w, http.StatusBadRequest, "Invalid request payload", nil, err)
		return
	}

	// Validate user data
	err = internal.ValidateProjectFields(&projectFields)
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
		Name:  projectFields.Name,
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

	projectFields = internal.ProjectInputs{}
	// Return success response
	sendJSONResponse(w, http.StatusCreated, "Project created successfully", project, nil)
	// TODO, Make the project unique on the user and the env type -> [test, dev..etc]
}

func (a *App) updateProjectHandler(w http.ResponseWriter, r *http.Request) {
	// Parse project ID from URL parameters
	vars := mux.Vars(r)
	projectIDStr := vars["id"]
	convertedProjectId, err := strconv.ParseInt(projectIDStr, 10, 64)
	if err != nil {
		sendJSONResponse(w, http.StatusBadRequest, "Cannot convert project id to number.", nil, err)
	}

	projectID := int(convertedProjectId)
	var updatedProject models.Project

	err = json.NewDecoder(r.Body).Decode(&updatedProject)
	if err != nil {
		sendJSONResponse(w, http.StatusBadRequest, "Invalid request body", nil, err)
		return
	}

	existingProject, err := a.DB.GetProjectByID(projectID)
	if err != nil {
		sendJSONResponse(w, http.StatusNotFound, fmt.Sprintf("Failed to retrieve project with id %s.", projectIDStr), nil, err)
		return
	}

	existingProject = updatedProject
	existingProject.ID = projectID

	err = a.DB.UpdateProject(existingProject)
	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, "Failed to update project", nil, err)
		return
	}
	sendJSONResponse(w, http.StatusOK, "Project updated successfully", existingProject, nil)
}
