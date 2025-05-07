package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	// import the generated protobuf code
)

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Post("/orders", createOrderHandler)
	router.Get("/orders/{order_id}", getOrderHandler)

	http.ListenAndServe(":8080", router)
}
