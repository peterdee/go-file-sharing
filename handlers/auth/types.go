package auth

type SetUpRequestPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInRequestPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
