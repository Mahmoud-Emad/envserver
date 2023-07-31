package app

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// deleteUserByIDHandler handles the HTTP request for deleting a user by ID.
// It expects the user ID to be present as a path parameter.
// If the user is successfully deleted, it returns a JSON response with status 204 (No Content).
// If the user is not found, it returns a JSON response with status 404 (Not Found).
// If there is an error during the deletion process, it returns an appropriate error response.
func (a *App) deleteUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from the path parameters
	vars := mux.Vars(r)
	if len(vars) == 0 {
		sendJSONResponse(w, http.StatusBadRequest, "Cannot get the user id.", nil, errors.New("User id should be provided."))
	}
	userIDStr := vars["id"]
	// Convert the user ID to uint
	u64, err := strconv.ParseUint(userIDStr, 10, 32)
	uID := uint(u64)

	// Check if the user exists
	_, err = a.DB.GetUserByID(uID)
	if err != nil {
		sendJSONResponse(w, http.StatusNotFound, "User not found", nil, err)
		return
	}

	// Delete the user from the database
	err = a.DB.DeleteUserByID(uID)
	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, "Failed to delete user", nil, err)
		return
	}

	// Return success response
	sendJSONResponse(w, http.StatusNoContent, "User deleted successfully", nil, nil)
}

// getUsersHandler handles the HTTP request for retrieving all users from the database.
// It returns a JSON response with status 200 (OK) containing an array of users.
// If the retrieval encounters an error, it returns an appropriate error response.
func (a *App) getUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := a.DB.GetUsers()
	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, "Failed to retrieve users", nil, err)
		return
	}

	sendJSONResponse(w, http.StatusOK, "Users found", users, nil)
}
