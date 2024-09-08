package controllers

import (
	"fmt"
	"net/http"

	"github.com/Dom-HTG/warp/middlewares"
	"github.com/Dom-HTG/warp/utils"
	"github.com/sirupsen/logrus"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	//This function handles the user's profile page.
	//Retrieve access token from request context.
	token, ok := r.Context().Value(middlewares.TokenKey).(string)
	if !ok {
		logrus.Fatal("Unable to retrieve access token from context")
	}

	userProfile, err := utils.GetUserProfile(token, r.Context())
	if err != nil {
		logrus.Fatal("Unable to retrieve access token from context: ", err)
	}

	logrus.Info("User profile retrieved")
	fmt.Print(userProfile)
}
