package models

type Comment struct {
	UserId  int    `json:"userId" db:"user_id"`
	PostId  int    `json:"postId" db:"post_id"`
	Content string `json:"content" db:"content"`
	BaseModel
}

type Reply struct {
	UserId    int    `json:"userId" db:"user_id"`
	CommentId int    `json:"commentId" db:"comment_id"`
	Content   string `json:"content" db:"content"`
	BaseModel
}

// type Person struct {
// 	Id int
// 	Fullname string `db:"full_name"`
// 	Country `db:"country"`
// }

// type Country struct {
// 	Id int
// 	Name string `db:"name"`
// }

// q = `select p.id, p.full_name, c.id as "country.id", c.name as "country.name" from person p inner join country c on p.country_id = c.id`

// var p []Person
// db.Select(&p, q)
