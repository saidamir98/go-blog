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

var CreatReply = func(w http.ResponseWriter, r *http.Request) {
	var reply models.Reply
	err := json.NewDecoder(r.Body).Decode(&reply)
	if err != nil {
		u.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	c := context.Get(r, "user")
	user := c.(*models.JwtCustomClaims)
	reply.UserId = user.Id

	q := `INSERT INTO replies (user_id, comment_id, content) values (:user_id, :comment_id, :content)
	RETURNING id, user_id, comment_id, content, created_at, updated_at`
	rows, err := app.DB.NamedQuery(q, reply)
	if err != nil {
		u.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	hasResult := false
	for rows.Next() {
		hasResult = true
		err := rows.StructScan(&reply)
		if err != nil {
			u.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
	}
	if !hasResult {
		u.RespondError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	u.RespondJSON(w, http.StatusOK, reply)
}

var UpdateReply = func(w http.ResponseWriter, r *http.Request) {
	var reply models.Reply
	err := json.NewDecoder(r.Body).Decode(&reply)
	if err != nil {
		u.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	c := context.Get(r, "user")
	user := c.(*models.JwtCustomClaims)
	reply.UserId = user.Id

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		u.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	reply.Id = id
	q := `
	UPDATE
		replies
	SET
		comment_id=:comment_id, content=:content
	WHERE 
		id=:id AND user_id=:user_id
	RETURNING
		id, user_id, comment_id, content, created_at, updated_at
	`
	rows, err := app.DB.NamedQuery(q, reply)
	if err != nil {
		u.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	hasResult := false
	for rows.Next() {
		hasResult = true
		err := rows.StructScan(&reply)
		if err != nil {
			u.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
	}
	if !hasResult {
		u.RespondError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	u.RespondJSON(w, http.StatusOK, reply)
}

var DeleteReply = func(w http.ResponseWriter, r *http.Request) {
	var dReply struct {
		Id     int `json:"id" db:"id"`
		UserId int `json:"userId" db:"user_id"`
	}

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		u.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	dReply.Id = id

	c := context.Get(r, "user")
	user := c.(*models.JwtCustomClaims)
	dReply.UserId = user.Id

	q := `
		DELETE 
		FROM 
		  replies c 
		WHERE
			c.id = :id AND c.user_id = :user_id
		`
	res, err := app.DB.NamedExec(q, dReply)
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

	u.RespondJSON(w, http.StatusOK, "REPLY has been deleted")
}
