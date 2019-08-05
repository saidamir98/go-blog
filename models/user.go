package models

import (
	"crypto/subtle"
	"encoding/hex"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	app "github.com/saidamir98/blog/app"
	"golang.org/x/crypto/argon2"
)

type BaseModel struct {
	Id        uint       `json:"id" db:"id"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}

type User struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
	Role     uint   `json:"role" db:"role"`
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
	Id   uint `json:"id"`
	Role uint `json:"role"`
	jwt.StandardClaims
}

func (u *User) GenerateUserJwt() (string, error) {
	claims := &JwtCustomClaims{
		u.Id,
		u.Role,
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

// func Insert(args ...string) {
// 	Insert("email", "test@gmil.com", "username", "farrukh")
// 	var a = "insert into tablename"
// 	var colNames string
// 	var colValues string
// 	for i, k := range args {
// 			if i % 2 != 0 {
// 				colNames += fmt.Sprintf("%s, ", args[k])
// 			} else {
// 				colValues += fmt.Sprintf("%s, ", args[k])
// 			}
// 	}
// 	colNames == fmt.Sprintf("(%s)", "id, username, ")
// 	a + colNames + "values" +  colValues
// 	fmt.Sprintf("insert into aaa (%s) values (%s)", args[0], args[1])
// }

// func (u *User) Save() (*User, error) {
// 	q := `INSERT INTO users (email, password, role, active)
// 	VALUES (:email, :password, :role, :active)`
// 	v := map[string]interface{}{
// 		"email":    u.Email,
// 		"password": u.Password,
// 		"role":     u.Role,
// 		"active":   u.Active,
// 	}
// 	rows, err := app.DB.NamedQuery(q, v)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if rows.Next() {
// 		err = rows.StructScan(&u)
// 	}
// 	if err != nil {
// 		return nil, err
// 	}
// 	return u, nil
// }

// func ListAllUsers() (users []User, err error) {
// 	err = app.DB.Select(&users, `SELECT * FROM users`)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return users, nil
// }
