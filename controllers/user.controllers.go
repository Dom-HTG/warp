package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

type AuthParams struct {
	ClientID     string
	ResponseType string
	RedirectURI  string
	State        string
	Scope        string
	ShowDialog   string
}

func getState() string {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	//Get authorization endpoint.
	baseURL := os.Getenv("BASE_URL")
	auth_endpoint := fmt.Sprintf("%s/authotize", baseURL)
	u, err := url.Parse(auth_endpoint)
	if err != nil {
		log.Fatal(err)
	}

	newAuthParams := &AuthParams{
		ClientID:     os.Getenv("CLIENT_ID"),
		ResponseType: "code",
		RedirectURI:  os.Getenv("REDIRECT_URI"),
		State:        getState(),
		Scope:        "user-top-read",
		ShowDialog:   "false",
	}

	//create query object and populate queries.
	query := u.Query()
	query.Set("client_id", newAuthParams.ClientID)
	query.Set("response_type", newAuthParams.ResponseType)
	query.Set("redirect_uri", newAuthParams.RedirectURI)
	query.Set("state", newAuthParams.State)
	query.Set("scope", newAuthParams.Scope)
	query.Set("show_dialog", newAuthParams.ShowDialog)

	//Append queries to parsed url.
	u.RawQuery = query.Encode()

	newURL := u.String()

	http.Redirect(w, r, newURL, http.StatusFound)
}
