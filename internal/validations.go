package internal

import (
	"fmt"
	"reflect"
)

// ValidateUser checks for the presence of required fields in the user struct.
func ValidateUserFields(user *SignUpInputs) error {
	t := reflect.TypeOf(*user)
	v := reflect.ValueOf(*user)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Name != "ProjectOwner" {
			value := v.Field(i).Interface()
			// Check if the field value is empty or zero
			if reflect.DeepEqual(value, reflect.Zero(field.Type).Interface()) {
				return fmt.Errorf("%s field is required", field.Name)
			}
		}
	}

	return nil
}

// ValidateProjectFields checks for the presence of required fields in the project struct.
func ValidateProjectFields(project *ProjectInputs) error {
	t := reflect.TypeOf(*project)
	v := reflect.ValueOf(*project)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()
		// Check if the field value is empty or zero
		if reflect.DeepEqual(value, reflect.Zero(field.Type).Interface()) {
			return fmt.Errorf("%s field is required", field.Name)
		}
	}

	return nil
}

// CheckPasswordHash compares the given plain-text password with the stored hashed password.
// It returns true if the password matches the hash, otherwise false.
func CheckPasswordHash(password string, hashedPassword []byte) bool {
	// Hash the user password
	mdPassHash := HashMD5(password)
	_, err := DecryptAES(hashedPassword, mdPassHash)
	return err == nil
}
