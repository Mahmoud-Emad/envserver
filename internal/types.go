package internal

// SignUpInputs struct for data needed when user create an account
type SignUpInputs struct {
	Name         string `json:"name" binding:"required" validate:"min=3,max=20"`
	Email        string `json:"email" binding:"required" validate:"mail"`
	Password     string `json:"password" binding:"required" validate:"password"`
	ProjectOwner bool   `json:"is_owner"`
}

// SigninInputs represents the input data for the user signin process.
type SigninInputs struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
