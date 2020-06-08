package model

type UserAuthenticationPayload struct {
	Profile      *User  `json:"profile"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}
