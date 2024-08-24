package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Dom-HTG/warp/controllers"
	"github.com/gorilla/mux"
)

var port string = os.Getenv("APP_PORT")

func main() {
	//Instantiate new router.
	r := mux.NewRouter()

	//Handle request multiplexing.
	r.HandleFunc("/signin", controllers.SignInHandler).Methods("GET")
	r.HandleFunc("/callback", controllers.CallbackHandler).Methods()

	//Run the server.
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("server is running on port %s\n", port)
}
