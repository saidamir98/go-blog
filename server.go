package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	app "github.com/saidamir98/blog/app"
	models "github.com/saidamir98/blog/models"
	routes "github.com/saidamir98/blog/routes"
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

	http.Handle("/", routes.Handlers())

	log.Printf("On address [%s] webServer is running...\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
