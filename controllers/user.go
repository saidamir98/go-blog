package controllers

import (
	"encoding/json"
	"net/http"

	app "github.com/saidamir98/blog/app"
	models "github.com/saidamir98/blog/models"
	u "github.com/saidamir98/blog/utils"
)

var Test = func(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	q := `SELECT * FROM users`
	err := app.DB.Select(&users, q)

	if err != nil {
		u.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	u.RespondJSON(w, http.StatusOK, users)
}

var Register = func(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		u.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	user.SetPassword(user.Password)
	user.RoleId = 1
	user.Active = true

	q := `
	INSERT INTO users 
		(username, email, password, role_id, active)
	VALUES
		(:username, :email, :password, :role_id, :active)
	`
	_, err = app.DB.NamedExec(q, user)
	if err != nil {
		u.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	q = `
	SELECT 
		*
	FROM
		users 
	WHERE
		username = $1
	LIMIT 
		1`
	err = app.DB.Get(&user, q, user.Username)
	if err != nil {
		u.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	var payload struct {
		Token string `json:"token"`
	}

	payload.Token, err = user.GenerateUserJwt()
	if err != nil {
		u.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	u.RespondJSON(w, http.StatusOK, payload)
}

var Login = func(w http.ResponseWriter, r *http.Request) {
	var userCredentials struct {
		Username string `json:"username" db:"username"`
		Password string `json:"password" db:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&userCredentials)
	if err != nil {
		u.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	var user models.User
	q := `
		SELECT 
			*
		FROM 
		  users u 
		WHERE
			u.username = $1 AND u.active
		LIMIT
			1`
	err = app.DB.Get(&user, q, userCredentials.Username)
	errorMessage := "Username or password incorrect"
	if err != nil {
		u.RespondError(w, http.StatusBadRequest, errorMessage)
		return
	}

	if ok := user.CheckPassword(userCredentials.Password); !ok {
		u.RespondError(w, http.StatusUnauthorized, errorMessage)
		return
	}

	var payload struct {
		Token string `json:"token"`
	}

	payload.Token, err = user.GenerateUserJwt()
	if err != nil {
		u.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	u.RespondJSON(w, http.StatusOK, payload)
}
