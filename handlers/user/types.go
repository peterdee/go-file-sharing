package user

type ChangePasswordRequestPayload struct {
	NewPassword string `json:"newPassword"`
	OldPassword string `json:"oldPassword"`
}
