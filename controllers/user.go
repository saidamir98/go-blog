package controllers

import (
	"encoding/json"
	"net/http"

	models "github.com/saidamir98/blog/models"
	u "github.com/saidamir98/blog/utils"
)

var CreateUser = func(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		u.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	token, err := user.GenerateUserJwt()
	if err != nil {
		u.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	u.RespondJSON(w, http.StatusOK, token)
}
