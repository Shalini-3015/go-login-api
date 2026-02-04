package auth

type RegisterUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterUserResponse struct {
	Message string `json:"message"`
}
