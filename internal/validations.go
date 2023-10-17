package internal

import (
	"fmt"
	"reflect"

	"golang.org/x/crypto/bcrypt"
)

type Validatable interface {
	Validate() error
}

func ValidateFields(v Validatable) error {
	rv := reflect.ValueOf(v).Elem()
	rt := rv.Type()

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		if field.Name != "ProjectOwner" {
			// ProjectOwner will be taken from the requested user
			value := rv.Field(i).Interface()

			// Check if the field value is empty or zero
			if reflect.DeepEqual(value, reflect.Zero(field.Type).Interface()) {
				return fmt.Errorf("%s field is required", field.Name)
			}

		}
	}

	return nil
}

// ValidateUser checks for the presence of required fields in the user struct.
func (s *SignUpInputs) Validate() error {
	return ValidateFields(s)
}

// ValidateProjectEnv checks for the presence of required fields in the env inputs struct.
func (e *EnvironmentKeyInputs) Validate() error {
	return ValidateFields(e)
}

// ValidateProjectFields checks for the presence of required fields in the project struct.
func (p *ProjectInputs) Validate() error {
	return ValidateFields(p)
}

// HashPassword hashes the given plain-text password using bcrypt.
// It returns the hashed password or an error if hashing fails.
func HashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashedPassword, nil
}

// CheckPasswordHash compares the given plain-text password with the stored hashed password.
// It returns true if the password matches the hash, otherwise false.
func CheckPasswordHash(password string, hashedPassword []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	return err == nil
}
