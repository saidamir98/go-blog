package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	app "github.com/saidamir98/blog/app"
	models "github.com/saidamir98/blog/models"
	u "github.com/saidamir98/blog/utils"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	app.Conf, err = godotenv.Read()
	if err != nil {
		log.Fatal("Error reading .env file")
	}
	models.InitDB()
}

func main() {
	addr := fmt.Sprintf("%s:%s", app.Conf["IP_ADDRESS"], app.Conf["PORT"])
	if app.Mux == nil {
		app.Mux = mux.NewRouter()
	}
	app.Mux.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	app.Mux.HandleFunc("/test", CreateUser).Methods("POST")
	// app.Mux.HandleFunc("/login", operators.Login).Methods("GET", "POST")
	// app.Mux.HandleFunc("/logout", operators.Logout).Methods("GET")

	log.Printf("On address [%s] webServer is running...\n", addr)
	err := http.ListenAndServe(addr, app.Mux)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

var CreateUser = func(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		u.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	// _, err = user.Save();
	// if err != nil {
	// 	u.RespondError(w, http.StatusBadRequest, err.Error())
	// 	return
	// }
	// var users []models.User
	// users, err = models.ListAllUsers()
	// if err != nil {
	// 	u.RespondError(w, http.StatusBadRequest, err.Error())
	// 	return
	// }
	u.RespondJSON(w, http.StatusOK, user)
}
