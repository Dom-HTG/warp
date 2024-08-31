package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Dom-HTG/warp/utils"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	//This function handles the user's profile page.
	//Retrieve access token from request context.
	accessToken, ok := r.Context().Value("accessToken").(string)
	if !ok {
		log.Fatal("Unable to retrieve access token from context.")
	}

	userProfile, err := utils.GetUserProfile(accessToken)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(userProfile)
}
