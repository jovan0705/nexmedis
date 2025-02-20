package Routers

import (
	"github.com/jovan0705/nexmedis/controllers"
	"github.com/jovan0705/nexmedis/middleware"

	"github.com/gorilla/mux"
)

func CartRouter(r *mux.Router) *mux.Router {
    cartController := &controllers.CartController{}

	route_cart := r.PathPrefix("/cart").Subrouter()
	route_cart.Use(middleware.IsAuthenticated)

    route_cart.HandleFunc("", cartController.AddToCart).Methods("POST")
	route_cart.HandleFunc("/{userId}", cartController.GetCart).Methods("GET")
	route_cart.HandleFunc("/{userId}", cartController.PayCartItems).Methods("PATCH")
	return r
}
