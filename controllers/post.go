package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

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

var GetPost = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		u.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	var post models.Post
	q := `
		SELECT 
			*
		FROM 
		  posts p 
		WHERE
			p.id = $1`
	err = app.DB.Get(&post, q, id)
	if err != nil {
		u.RespondError(w, http.StatusBadRequest, err.Error())
		return
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
	post.Id = (uint)(id)
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
		Id     int  `json:"id" db:"id"`
		UserId uint `json:"userId" db:"user_id"`
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
