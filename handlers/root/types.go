package root

type CreateUserRequestPayload struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}
