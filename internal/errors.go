package internal

import (
	"errors"
	"fmt"
)

var (
	cantLoadConfigFileError   = errors.New("failed to open config file, Please make sure that you have a config file called config.toml in your main root, please see the ./config.toml.template")
	cantDecodeConfigError     = errors.New("failed to decode config from reader")
	InternalServerError       = errors.New("something went wrong")
	UserEmailNotUniqueError   = errors.New("the user email field must be unique")
	UserIdNotProvidedError    = errors.New("user id should be provided")
	ProjectIdNotProvidedError = errors.New("project id should be provided")
)

func missingKeyError(keyName string) error {
	return fmt.Errorf("the %s key is missing in the config file", keyName)
}
