package internal

import (
	"errors"
	"strings"
)

var (
	serverHostKeyError = errors.New("The server Host key is missing in the config file.")
	serverPortKeyError = errors.New("The server port key is missing in the config file.")

	databaseHostKeyError     = errors.New("The database Host key is missing in the config file.")
	databasePortKeyError     = errors.New("The database port key is missing in the config file.")
	databaseNameKeyError     = errors.New("The database name key is missing in the config file.")
	databaseUserKeyError     = errors.New("The database user key is missing in the config file.")
	databasePasswordKeyError = errors.New("The database password key is missing in the config file.")
	InternalServerError      = errors.New("Something went wrong")
	UserEmailNotUniqueError  = errors.New("The user email field must be unique.")
)

// Check all required fields then return an error if there.
func (c *Configuration) validateConfiguration() error {
	if strings.TrimSpace(c.Server.Host) == "" {
		return serverHostKeyError
	}

	if c.Server.Port == 0 {
		// That's mean the port key is not there.
		return serverPortKeyError
	}

	if strings.TrimSpace(c.Database.Host) == "" {
		return databaseHostKeyError
	}

	if strings.TrimSpace(c.Database.Name) == "" {
		return databaseNameKeyError
	}

	if strings.TrimSpace(c.Database.User) == "" {
		return databaseUserKeyError
	}

	if strings.TrimSpace(c.Database.Password) == "" {
		return databasePasswordKeyError
	}

	if c.Database.Port == 0 {
		return databasePortKeyError
	}

	return nil
}
