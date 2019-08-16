package middlewares

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/context"

	jwt "github.com/dgrijalva/jwt-go"

	app "github.com/saidamir98/go-blog/app"
	models "github.com/saidamir98/go-blog/models"
	u "github.com/saidamir98/go-blog/utils"
)

var JwtVerify = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		notAuth := []string{"/register", "/login"}
		requestPath := r.URL.Path
		for _, value := range notAuth {
			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		tokenHeader := r.Header.Get("Authorization")
		if tokenHeader == "" {
			u.RespondError(w, http.StatusForbidden, "Missing auth token")
			return
		}
		splitted := strings.Split(tokenHeader, " ")
		if len(splitted) != 2 {
			u.RespondError(w, http.StatusForbidden, "Invalid/Malformed auth token")
			return
		}

		var tk models.JwtCustomClaims
		token, err := jwt.ParseWithClaims(splitted[1], &tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(app.Conf["JWT_SECRET"]), nil
		})
		if err != nil {
			u.RespondError(w, http.StatusForbidden, "Malformed authentication token")
			return
		}
		if !token.Valid {
			u.RespondError(w, http.StatusForbidden, "Token is not valid")
			return
		}

		log.Printf("User id %v", tk.Id)

		context.Set(r, "user", &tk)
		next.ServeHTTP(w, r)
	})
}
