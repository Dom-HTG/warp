package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/Dom-HTG/warp/models"
	"github.com/Dom-HTG/warp/utils"
	"gorm.io/gorm"
)

type Handlers interface {
	SignInHandler(w http.ResponseWriter, r *http.Request)
	CallbackHandler(w http.ResponseWriter, r *http.Request)
}

type repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *repo {
	return &repo{
		db: db,
	}
}

var globalStateID uint

func (rp repo) SignInHandler(w http.ResponseWriter, r *http.Request) {
	//Get authorization endpoint.
	baseURL := os.Getenv("BASE_URL")
	u, err := url.Parse(baseURL)
	if err != nil {
		log.Fatal(err)
	}

	stateData := utils.GenerateState()

	newAuthParams := &models.AuthParams{
		ClientID:     os.Getenv("CLIENT_ID"),
		ResponseType: "code",
		RedirectURI:  os.Getenv("REDIRECT_URI"),
		State:        stateData,
		Scope:        "user-top-read user-read-private user-read-email",
		//user-top-read allows access to the users' top tracks and top artists.
		//user-read-private allows access to the users' private data.
		//user-read-email allows access to the users' email address.
		ShowDialog: "false",
	}

	//Store state to the database.
	state := &models.User{
		StateValue: stateData,
	}

	tx := rp.db.Create(&state)
	if tx.Error != nil {
		log.Fatalf("unable to save state to DB: %v", tx.Error)
	}

	//store state ID in global variable.
	globalStateID = state.ID

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

func (rp repo) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	//This callback retrieves the code(authorization code) and state sent back by the spotify service.
	//Create a query object on the request object.
	requery := r.URL.Query()
	authCode := requery.Get("code")
	state := requery.Get("state")

	//get state ID from globalStateID variable.
	id := globalStateID

	//get state values from database.
	DBstate, err := utils.GetStateDB(rp.db, id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("state fetched from database. \n")

	//compare state values.
	if state != DBstate {
		log.Fatal("state mismatched")
	}
	fmt.Print("state MATCHED. \n")

	//Exchange authorization code for access token and refresh token.
	tokenPayload, err1 := utils.GetAccessToken(authCode)
	if err1 != nil {
		log.Fatal(err1)
	}
	if tokenPayload == nil {
		fmt.Printf("no access and refresh tokens. \n")
	}

	//token data to be committed to context.
	tokenContext := &models.TokenContext{
		AccessToken:  tokenPayload.AccessToken,
		RefreshToken: tokenPayload.RefreshToken,
	}

	//commit token to context.
	ctx := context.WithValue(r.Context(), "access_token", tokenContext)
	ProfileHandler(w, r.WithContext(ctx))

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

func (rp repo) HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("welcome to warp home"))
}
