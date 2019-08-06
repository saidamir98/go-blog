package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/context"
	app "github.com/saidamir98/blog/app"
	models "github.com/saidamir98/blog/models"
	u "github.com/saidamir98/blog/utils"
)

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
	for rows.Next() {
		err := rows.StructScan(&post)
		if err != nil {
			log.Fatalln(err)
		}
	}
	if err != nil {
		u.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	u.RespondJSON(w, http.StatusOK, post)
}
