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

var CreatComment = func(w http.ResponseWriter, r *http.Request) {
	var comment models.Comment
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		u.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	c := context.Get(r, "user")
	user := c.(*models.JwtCustomClaims)
	comment.UserId = user.Id

	q := `INSERT INTO comments (user_id, post_id, content) values (:user_id, :post_id, :content)
	RETURNING id, user_id, post_id, content, created_at, updated_at`
	rows, err := app.DB.NamedQuery(q, comment)
	if err != nil {
		u.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	hasResult := false
	for rows.Next() {
		hasResult = true
		err := rows.StructScan(&comment)
		if err != nil {
			u.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
	}
	if !hasResult {
		u.RespondError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	u.RespondJSON(w, http.StatusOK, comment)
}

var UpdateComment = func(w http.ResponseWriter, r *http.Request) {
	var comment models.Comment
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		u.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	c := context.Get(r, "user")
	user := c.(*models.JwtCustomClaims)
	comment.UserId = user.Id

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		u.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	comment.Id = id
	q := `
	UPDATE
		comments
	SET
		post_id=:post_id, content=:content
	WHERE 
		id=:id AND user_id=:user_id
	RETURNING
		id, user_id, post_id, content, created_at, updated_at
	`
	rows, err := app.DB.NamedQuery(q, comment)
	if err != nil {
		u.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	hasResult := false
	for rows.Next() {
		hasResult = true
		err := rows.StructScan(&comment)
		if err != nil {
			u.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
	}
	if !hasResult {
		u.RespondError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	u.RespondJSON(w, http.StatusOK, comment)
}

var DeleteComment = func(w http.ResponseWriter, r *http.Request) {
	var dComment struct {
		Id     int `json:"id" db:"id"`
		UserId int `json:"userId" db:"user_id"`
	}

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		u.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	dComment.Id = id

	c := context.Get(r, "user")
	user := c.(*models.JwtCustomClaims)
	dComment.UserId = user.Id

	q := `
		DELETE 
		FROM 
		  comments c 
		WHERE
			c.id = :id AND c.user_id = :user_id
		`
	res, err := app.DB.NamedExec(q, dComment)
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

	u.RespondJSON(w, http.StatusOK, "COMMENT has been deleted")
}
