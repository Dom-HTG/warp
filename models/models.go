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

type TokenContext struct {
	AccessToken  string
	RefreshToken string
}

type UserProfile struct {
	ID        string     `json:"id"`
	Name      string     `json:"display_name"`
	Email     string     `json:"email"`
	Href      string     `json:"href"`
	URI       string     `json:"uri"`
	Followers []follower `json:"followers"`
	Images    []image    `json:"images"`
}

type follower struct {
	Href  string `json:"href"`
	Total int    `json:"total"`
}

type image struct {
	ImageURL string `json:"url"`
	Height   int    `json:"height"`
	Width    int    `json:"width"`
}
