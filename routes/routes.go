package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/saidamir98/blog/controllers"
	"github.com/saidamir98/blog/middlewares"
)

func Handlers() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.Use(CommonMiddleware)
	// r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	r.HandleFunc("/", controllers.Test).Methods("GET")
	// r.HandleFunc("/api", controllers.TestAPI).Methods("GET")
	r.HandleFunc("/register", controllers.Register).Methods("POST")
	r.HandleFunc("/login", controllers.Login).Methods("POST")

	s := r.PathPrefix("/auth").Subrouter()
	s.Use(middlewares.JwtVerify)
	// s.HandleFunc("/test", controllers.Test).Methods("POST")

	s.HandleFunc("/my-posts", controllers.ListUserPosts).Methods("GET")
	s.HandleFunc("/posts", controllers.CreatPost).Methods("POST")
	s.HandleFunc("/posts", controllers.ListPosts).Methods("GET")
	s.HandleFunc("/posts/{id}", controllers.GetPost).Methods("GET")
	s.HandleFunc("/posts/{id}", controllers.UpdatePost).Methods("PUT")
	s.HandleFunc("/posts/{id}", controllers.DeletePost).Methods("DELETE")

	s.HandleFunc("/comments", controllers.CreatComment).Methods("POST")
	s.HandleFunc("/comments/{id}", controllers.UpdateComment).Methods("PUT")
	s.HandleFunc("/comments/{id}", controllers.DeleteComment).Methods("DELETE")

	s.HandleFunc("/replies", controllers.CreatReply).Methods("POST")
	s.HandleFunc("/replies/{id}", controllers.UpdateReply).Methods("PUT")
	s.HandleFunc("/replies/{id}", controllers.DeleteReply).Methods("DELETE")
	return r
}

// CommonMiddleware --Set content-type
func CommonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
		next.ServeHTTP(w, r)
	})
}
