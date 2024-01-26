package internal

// SignUpInputs struct for data needed when user create an account
type SignUpInputs struct {
	FirstName    string `json:"first_name" binding:"required" validate:"min=3,max=20"`
	LastName     string `json:"last_name" binding:"required" validate:"min=3,max=20"`
	Email        string `json:"email" binding:"required" validate:"mail"`
	Password     string `json:"password" binding:"required" validate:"password"`
	ProjectOwner bool   `json:"is_owner"`
}

// SigninInputs represents the input data for the user signin process.
type SigninInputs struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ProjectInputs represents the input data for the create project process.
type ProjectInputs struct {
	Name string `json:"name"`
}

type EnvironmentKeyInputs struct {
	Key   string
	Value string
}
