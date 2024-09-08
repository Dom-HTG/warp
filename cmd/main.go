package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Dom-HTG/warp/controllers"
	"github.com/Dom-HTG/warp/utils"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	//load .env file.
	err0 := godotenv.Load()
	if err0 != nil {
		log.Fatal(err0)
	}

	//Initialize logging.
	utils.InitLogger()

	//Instantiate new router.
	r := mux.NewRouter()

	//Instantiate new DB connection.
	db, err := utils.InitDB()
	if err != nil {
		logrus.Errorf("Error initializing DB: %v", err)
	}
	logrus.Info("Database initialized and connected successfully")
	logrus.Info("Model migration success")

	//depependency injection.
	controller := controllers.NewRepo(db)

	//Apply middleware.

	//Auth routes.
	r.HandleFunc("/signin", controller.SignInHandler).Methods("GET")
	r.HandleFunc("/callback", controller.CallbackHandler).Methods("GET")

	r.HandleFunc("/home", controller.HomeHandler).Methods("GET")

	//Data retrieval routes.
	user := r.PathPrefix("/user").Subrouter()
	user.HandleFunc("/profile", controllers.ProfileHandler).Methods("GET")

	//Run the server.
	port := os.Getenv("APP_PORT")
	logrus.Info("Server is running on port: ", port)
	err1 := http.ListenAndServe(fmt.Sprintf(":%s", port), r)
	if err1 != nil {
		logrus.Errorf("Error starting server: ", err1)
	}

}
