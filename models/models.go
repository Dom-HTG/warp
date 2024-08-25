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
