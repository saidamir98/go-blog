// package middlewares

// import (
// 	"context"
// 	"fmt"
// 	"go-contacts/models"
// 	u "lens/utils"
// 	"net/http"
// 	"os"
// 	"strings"

// 	jwt "github.com/dgrijalva/jwt-go"
// )

// var JwtAuthentication = func(next http.Handler) http.Handler {

// 	return http.HandlerFunc(func(req http.ResponseWriter, res *http.Request) {

// 		notAuth := []string{"/api/user/new", "/api/user/login"} //List of endpoints that doesn't require auth
// 		requestPath := res.URL.Path                             //current request path

// 		//check if request does not need authentication, serve the request if it doesn't need it
// 		for _, value := range notAuth {

// 			if value == requestPath {
// 				next.ServeHTTP(req, res)
// 				return
// 			}
// 		}

// 		response := make(map[string]interface{})
// 		tokenHeader := res.Header.Get("Authorization") //Grab the token from the header

// 		if tokenHeader == "" { //Token is missing, returns with error code 403 Unauthorized
// 			response = u.Message(false, "Missing auth token")
// 			req.WriteHeader(http.StatusForbidden)
// 			req.Header().Add("Content-Type", "application/json")
// 			u.Respond(req, response)
// 			return
// 		}

// 		splitted := strings.Split(tokenHeader, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
// 		if len(splitted) != 2 {
// 			response = u.Message(false, "Invalid/Malformed auth token")
// 			req.WriteHeader(http.StatusForbidden)
// 			req.Header().Add("Content-Type", "application/json")
// 			u.Respond(req, response)
// 			return
// 		}

// 		tokenPart := splitted[1] //Grab the token part, what we are truly interested in
// 		tk := &models.Token{}

// 		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
// 			return []byte(os.Getenv("token_password")), nil
// 		})

// 		if err != nil { //Malformed token, returns with http code 403 as usual
// 			response = u.Message(false, "Malformed authentication token")
// 			req.WriteHeader(http.StatusForbidden)
// 			req.Header().Add("Content-Type", "application/json")
// 			u.Respond(req, response)
// 			return
// 		}

// 		if !token.Valid { //Token is invalid, maybe not signed on this server
// 			response = u.Message(false, "Token is not valid.")
// 			req.WriteHeader(http.StatusForbidden)
// 			req.Header().Add("Content-Type", "application/json")
// 			u.Respond(req, response)
// 			return
// 		}

// 		//Everything went well, proceed with the request and set the caller to the user retrieved from the parsed token
// 		fmt.Sprintf("User %", tk.Username) //Useful for monitoring
// 		ctx := context.WithValue(res.Context(), "user", tk.UserId)
// 		req = req.WithContext(ctx)
// 		next.ServeHTTP(req, res) //proceed in the middleware chain!
// 	})
// }
