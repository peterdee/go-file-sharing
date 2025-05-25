package root

type CreateUserRequestPayload struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

type UpdateUserRequestPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
