package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
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
	port := os.Getenv("PORT")

	http.Handle("/", routes.Handlers())

	log.Printf("On port [%s] webServer is running...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS()(r)))
}
