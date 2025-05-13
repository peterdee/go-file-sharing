package auth

type setUpRequestPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type signInRequestPayload struct {
	setUpRequestPayload
}
