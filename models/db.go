package models

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	app "github.com/saidamir98/blog/app"
)

func InitDB() {
	var (
		err error
	)
	dbUri := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		app.Conf["DB_HOST"], app.Conf["DB_PORT"], app.Conf["DB_NAME"], app.Conf["DB_USERNAME"], app.Conf["DB_PASSWORD"])

	app.DB, err = sqlx.Connect("postgres", dbUri)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err = app.DB.Ping(); err != nil {
		log.Fatalf("%+v", err)
	}

	log.Println("connected db...")

	UserSchema := `
	DROP TABLE IF EXISTS users;
	CREATE TABLE users(
		id serial PRIMARY KEY,
		email VARCHAR (255) UNIQUE NOT NULL,
		password VARCHAR (50) NOT NULL,
		role integer,
		active BOOLEAN NOT NULL DEFAULT FALSE,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	`
	_, err = app.DB.Exec(UserSchema)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}
