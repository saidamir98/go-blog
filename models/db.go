package models

import (
	"log"
	"os"
	"time"

	// _ "github.com/jackc/pgx"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	app "github.com/saidamir98/blog/app"
)

type BaseModel struct {
	Id        int        `json:"id" db:"id"`
	CreatedAt *time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt *time.Time `json:"updatedAt" db:"updated_at"`
}

func InitDB() {
	var (
		err error
	)

	dbUri := os.Getenv("DATABASE_URL")

	app.DB, err = sqlx.Connect("postgres", dbUri)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err = app.DB.Ping(); err != nil {
		log.Fatalf("%+v", err)
	}

	log.Println("connected db...")

	// _, err = app.DB.Exec(Schemas)
	// if err != nil {
	// 	log.Fatalf("%+v", err)
	// }
	log.Println(Schemas)
}

var Schemas = `
	DROP TABLE IF EXISTS users;
	CREATE TABLE users(
		id serial PRIMARY KEY,
		username VARCHAR (255) UNIQUE NOT NULL,
		email VARCHAR (255) UNIQUE NOT NULL,
		password VARCHAR (255) NOT NULL,
		role_id INTEGER,
		active BOOLEAN NOT NULL DEFAULT FALSE,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	DROP TABLE IF EXISTS posts;
	CREATE TABLE posts(
		id serial PRIMARY KEY,
		title VARCHAR(50) UNIQUE NOT NULL,
		content TEXT UNIQUE NOT NULL,
		image VARCHAR(255),
		user_id INTEGER,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	DROP TABLE IF EXISTS comments;
	CREATE TABLE comments(
		id serial PRIMARY KEY,
		content TEXT UNIQUE NOT NULL,
		user_id INTEGER,
		post_id INTEGER,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	DROP TABLE IF EXISTS replies;
	CREATE TABLE replies(
		id serial PRIMARY KEY,
		content TEXT UNIQUE NOT NULL,
		user_id INTEGER,
		comment_id INTEGER,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	`
