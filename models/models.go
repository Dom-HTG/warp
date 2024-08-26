package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	StateValue string
}

type AuthParams struct {
	ClientID     string
	ResponseType string
	RedirectURI  string
	State        string
	Scope        string
	ShowDialog   string
}

type AccessTokenPayload struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}
