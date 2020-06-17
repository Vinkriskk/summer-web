package main

import (
	"fmt"
	"log"
	"net/http"

	delivery "summer-web/delivery/http"
	"summer-web/delivery/middleware"

	"github.com/gorilla/mux"
)

// 	set SECRET_JWT_KEY=super_secret_key
// 	set DB_CONNECTION_STRING=host=localhost port=5432 user=postgres dbname=summer_web_development password=password sslmode=disable

func main() {
	// initializeEnv()
	router := mux.NewRouter()

	var httpMiddleware middleware.Middleware = middleware.NewMiddleware()
	var postDelivery delivery.PostDelivery = delivery.NewPostDelivery()
	var userDelivery delivery.UserDelivery = delivery.NewUserDelivery()

	const port string = ":8000"

	router.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(resp, "Up and running...")
	})

	router.HandleFunc("/sign_up", userDelivery.AddUser).Methods("POST")
	router.HandleFunc("/login", userDelivery.Login).Methods("POST")

	router.Handle("/browse", httpMiddleware.IsAuthorized(postDelivery.GetPosts)).Methods("GET")
	router.Handle("/posts", httpMiddleware.IsAuthorized(postDelivery.AddPost)).Methods("POST")

	router.Handle("/users/{id}", httpMiddleware.IsAuthorized(userDelivery.GetUserByID)).Methods("GET")
	router.Handle("/users/update", httpMiddleware.IsAuthorized(userDelivery.UpdateUser)).Methods("PATCH")

	log.Println("Server is listening on port", port)
	log.Fatalln(http.ListenAndServe(port, router))
}
