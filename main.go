package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/jovan0705/nexmedis/controllers"
	"github.com/jovan0705/nexmedis/router"
)

func main() {
    userController := &controllers.UserController{}

    r := mux.NewRouter()

    r.HandleFunc("/api/register", userController.RegisterUser).Methods("POST")
    r.HandleFunc("/api/login", userController.LoginUser).Methods("POST")
	route := r.PathPrefix("/api").Subrouter()

	r.Handle("/", Routers.CartRouter(route))
	r.Handle("/", Routers.ProductRouter(route))

	log.Println("Server is starting on port 8080...")
    if err := http.ListenAndServe(":8080", r); err != nil {
        log.Fatal("Error starting server: ", err)
    }
}

