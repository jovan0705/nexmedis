package Routers

import (
	"github.com/jovan0705/nexmedis/controllers"
	middleware "github.com/jovan0705/nexmedis/middleware"

	"github.com/gorilla/mux"
)

func ProductRouter(r *mux.Router) *mux.Router {
	productController := &controllers.ProductController{}

	route_product := r.PathPrefix("/products").Subrouter()
	route_product.Use(middleware.IsAuthenticated)

	route_product.HandleFunc("/", productController.GetAllProducts).Methods("GET")
    route_product.HandleFunc("/{id}", productController.GetProduct).Methods("GET")
	return r
}
