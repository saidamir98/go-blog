package models

import (
	"crypto/subtle"
	"encoding/hex"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	app "github.com/saidamir98/blog/app"
	"golang.org/x/crypto/argon2"
)

type User struct {
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
	RoleId   int    `json:"roleId" db:"role_id"`
	Active   bool   `json:"active" db:"active"`
	BaseModel
}

func (u *User) SetPassword(password string) {
	key := argon2.Key([]byte(password), []byte(app.Conf["PASSWORD_SALT"]), 3, 32*1024, 4, 32)
	u.Password = hex.EncodeToString(key)
}

func (u *User) CheckPassword(password string) bool {
	key := argon2.Key([]byte(password), []byte(app.Conf["PASSWORD_SALT"]), 3, 32*1024, 4, 32)
	hashedPassword := hex.EncodeToString(key)
	if subtle.ConstantTimeCompare([]byte(u.Password), []byte(hashedPassword)) == 1 {
		return true
	}
	return false
}

type JwtCustomClaims struct {
	Id     int `json:"id"`
	RoleId int `json:"roleId"`
	jwt.StandardClaims
}

func (u *User) GenerateUserJwt() (string, error) {
	claims := &JwtCustomClaims{
		u.Id,
		u.RoleId,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(app.Conf["JWT_SECRET"]))
	if err != nil {
		return "", err
	}

	return t, nil
}
