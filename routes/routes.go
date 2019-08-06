package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/saidamir98/blog/controllers"
	"github.com/saidamir98/blog/middlewares"
)

func Handlers() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// r.HandleFunc("/", controllers.TestAPI).Methods("GET")
	// r.HandleFunc("/api", controllers.TestAPI).Methods("GET")
	r.HandleFunc("/register", controllers.Register).Methods("POST")
	// r.HandleFunc("/login", controllers.Login).Methods("POST")

	s := r.PathPrefix("/auth").Subrouter()
	s.Use(middlewares.JwtVerify)
	s.HandleFunc("/test", controllers.Test).Methods("POST")
	// s.HandleFunc("/user", controllers.FetchUsers).Methods("GET")
	// s.HandleFunc("/user/{id}", controllers.GetUser).Methods("GET")
	// s.HandleFunc("/user/{id}", controllers.UpdateUser).Methods("PUT")
	// s.HandleFunc("/user/{id}", controllers.DeleteUser).Methods("DELETE")
	return r
}
