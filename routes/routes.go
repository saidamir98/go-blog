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
	r.HandleFunc("/login", controllers.Login).Methods("POST")

	s := r.PathPrefix("/auth").Subrouter()
	s.Use(middlewares.JwtVerify)
	// s.HandleFunc("/test", controllers.Test).Methods("POST")
	s.HandleFunc("/posts", controllers.CreatPost).Methods("POST")
	s.HandleFunc("/posts", controllers.ListPosts).Methods("GET")
	s.HandleFunc("/posts/{id}", controllers.GetPost).Methods("GET")
	s.HandleFunc("/posts/{id}", controllers.UpdatePost).Methods("PUT")
	s.HandleFunc("/posts/{id}", controllers.DeletePost).Methods("DELETE")
	return r
}
