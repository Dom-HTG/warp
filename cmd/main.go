package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Dom-HTG/warp/controllers"
	"github.com/Dom-HTG/warp/middlewares"
	"github.com/Dom-HTG/warp/utils"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	//load .env file.
	err0 := godotenv.Load()
	if err0 != nil {
		log.Fatal(err0)
	}

	//Instantiate new router.
	r := mux.NewRouter()

	//Instantiate new DB connection.
	db, err := utils.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Database initialized and connected. \n")
	fmt.Print("model migration success. \n")

	//depependency injection.
	controller := controllers.NewRepo(db)

	//Auth routes.
	r.HandleFunc("/signin", controller.SignInHandler).Methods("GET")
	r.HandleFunc("/callback", controller.CallbackHandler).Methods("GET")

	r.HandleFunc("/home", controller.HomeHandler).Methods("GET")

	//Data retrieval routes.
	protected := r.PathPrefix("/home").Subrouter()
	protected.Use(middlewares.AddTokenToContext)
	protected.HandleFunc("/profile", controllers.ProfileHandler).Methods("GET")

	//Spotify Query routes.
	r.HandleFunc("/home/user-profile", controllers.ProfileHandler).Methods("GET")

	//Run the server.
	port := os.Getenv("APP_PORT")
	fmt.Printf("server is running on port %s \n", port)
	err1 := http.ListenAndServe(fmt.Sprintf(":%s", port), r)
	if err1 != nil {
		log.Fatal(err1)
	}

}
