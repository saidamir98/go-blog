package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	app "github.com/saidamir98/blog/app"
	models "github.com/saidamir98/blog/models"
	u "github.com/saidamir98/blog/utils"
)

var ListUserPosts = func(w http.ResponseWriter, r *http.Request) {
	c := context.Get(r, "user")
	user := c.(*models.JwtCustomClaims)

	var posts []models.Post
	q := `SELECT * FROM posts WHERE user_id=$1`
	err := app.DB.Select(&posts, q, user.Id)
	if err != nil {
		u.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if len(posts) <= 0 {
		u.RespondJSON(w, http.StatusOK, "No post found")
		return
	}

	u.RespondJSON(w, http.StatusOK, posts)
}

var CreatPost = func(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		u.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	c := context.Get(r, "user")
	user := c.(*models.JwtCustomClaims)
	post.UserId = user.Id

	q := `INSERT INTO posts (user_id, title, content, image) values (:user_id, :title, :content, :image)
	RETURNING id, user_id, title, content, image, created_at, updated_at`
	rows, err := app.DB.NamedQuery(q, post)
	if err != nil {
		u.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	hasResult := false
	for rows.Next() {
		hasResult = true
		err := rows.StructScan(&post)
		if err != nil {
			u.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
	}
	if !hasResult {
		u.RespondError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	u.RespondJSON(w, http.StatusOK, post)
}

var ListPosts = func(w http.ResponseWriter, r *http.Request) {
	var posts []models.Post
	q := `SELECT * FROM posts`
	err := app.DB.Select(&posts, q)
	if err != nil {
		u.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if len(posts) <= 0 {
		u.RespondJSON(w, http.StatusOK, "No post found")
		return
	}

	u.RespondJSON(w, http.StatusOK, posts)
}

type User struct {
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email"`
	Id       int    `json:"id" db:"id"`
}
type Reply struct {
	User      `json:"user" db:"user"`
	Content   string     `json:"content" db:"content"`
	Id        int        `json:"id" db:"id"`
	CreatedAt *time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt *time.Time `json:"updatedAt" db:"updated_at"`
}
type Comment struct {
	User      `json:"user" db:"user"`
	Content   string     `json:"content" db:"content"`
	Replies   []Reply    `json:"replies" db:"replies"`
	Id        int        `json:"id" db:"id"`
	CreatedAt *time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt *time.Time `json:"updatedAt" db:"updated_at"`
}
type Post struct {
	User      `json:"user" db:"user"`
	Title     string     `json:"title" db:"title"`
	Content   string     `json:"content" db:"content"`
	Image     string     `json:"image" db:"image"`
	Comments  []Comment  `json:"comments" db:"comments"`
	Id        int        `json:"id" db:"id"`
	CreatedAt *time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt *time.Time `json:"updatedAt" db:"updated_at"`
}
var GetPost = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		u.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	var post Post
	q := `
		SELECT 
			p.title, p.content, p.image, p.id, p.created_at, p.updated_at,
			u.username as "user.username",
			u.email as "user.email",
			u.id as "user.id"
		FROM 
			posts p
		INNER JOIN
			users u
		ON
			p.user_id = u.id 
		WHERE
			p.id = $1`
	err = app.DB.Get(&post, q, id)
	if err != nil {
		u.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	q = `SELECT 
			c.content, c.id, c.created_at, c.updated_at,
			u.username as "user.username",
			u.email as "user.email",
			u.id as "user.id"
		FROM 
			comments c
		INNER JOIN
			users u
		ON
			c.user_id = u.id 
		WHERE
			c.post_id = $1`
	err = app.DB.Select(&post.Comments, q, post.Id)
	if err != nil {
		u.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	for i := range post.Comments {
		q = `SELECT 
				r.content, r.id, r.created_at, r.updated_at,
				u.username as "user.username",
				u.email as "user.email",
				u.id as "user.id"
			FROM 
				replies r
			INNER JOIN
				users u
			ON
				r.user_id = u.id 
			WHERE
				r.comment_id = $1`
		err = app.DB.Select(&post.Comments[i].Replies, q, post.Comments[i].Id)
		if err != nil {
			u.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	u.RespondJSON(w, http.StatusOK, post)
}

var UpdatePost = func(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		u.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	c := context.Get(r, "user")
	user := c.(*models.JwtCustomClaims)
	post.UserId = user.Id

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		u.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	post.Id = id
	q := `
	UPDATE
		posts
	SET
		title=:title, content=:content, image=:image
	WHERE 
		id=:id AND user_id=:user_id
	RETURNING
		id, user_id, title, content, image, created_at, updated_at
	`
	rows, err := app.DB.NamedQuery(q, post)
	if err != nil {
		u.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	hasResult := false
	for rows.Next() {
		hasResult = true
		err := rows.StructScan(&post)
		if err != nil {
			u.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
	}
	if !hasResult {
		u.RespondError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	u.RespondJSON(w, http.StatusOK, post)
}

var DeletePost = func(w http.ResponseWriter, r *http.Request) {
	var dPost struct {
		Id     int `json:"id" db:"id"`
		UserId int `json:"userId" db:"user_id"`
	}

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		u.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	dPost.Id = id

	c := context.Get(r, "user")
	user := c.(*models.JwtCustomClaims)
	dPost.UserId = user.Id

	q := `
		DELETE 
		FROM 
		  posts p 
		WHERE
			p.id = :id AND p.user_id = :user_id
		`
	res, err := app.DB.NamedExec(q, dPost)
	if err != nil {
		u.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	count, err := res.RowsAffected()
	if err != nil {
		u.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	if count <= 0 {
		u.RespondError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	u.RespondJSON(w, http.StatusOK, "POST has been deleted")
}
