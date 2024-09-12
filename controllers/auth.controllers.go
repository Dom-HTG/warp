package controllers

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/Dom-HTG/warp/middlewares"
	"github.com/Dom-HTG/warp/models"
	"github.com/Dom-HTG/warp/utils"
	"github.com/sirupsen/logrus"
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
		logrus.Errorf("Error parsing auth base URL: %v", err)
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
		logrus.Errorf("Error storing state to databse: %v", tx.Error)
	}
	logrus.Info("state saved to database")

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
		logrus.Errorf("Error retrieving state from database: %v", err)
	}
	fmt.Printf("state fetched from database. \n")

	//compare state values.
	if state != DBstate {
		logrus.Fatal("State mismatch")
	}
	logrus.Info("State values matched")

	//Exchange authorization code for access token and refresh token.
	tokenPayload, err1 := utils.GetAccessToken(authCode, r.Context())
	if err1 != nil {
		logrus.Errorf("Error getting access token: %v", err1)
	}
	logrus.Info("Access token token Obtained")

	//token data to be committed to context.
	if tokenPayload == nil {
		logrus.Fatal("Access token payload is empty")
	}

	accessToken := tokenPayload.AccessToken

	//commit token to context.
	ctx := context.WithValue(r.Context(), utils.Token, accessToken)

	// ProfileHandler(w, r.WithContext(ctx))
	cx := r.WithContext(ctx)

	//passing updated context to the context middleware.
	middlewares.AddTokenToContext(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.Info("updated context passed to context middleware")
	})).ServeHTTP(w, cx)

	http.Redirect(w, cx, "/home", http.StatusFound)
}

func (rp repo) HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("welcome to warp home. You have authorized access to your spotify listening data"))
}
