package models

type Post struct {
	UserId  int    `json:"userId" db:"user_id"`
	Title   string `json:"title" db:"title"`
	Content string `json:"content" db:"content"`
	Image   string `json:"image" db:"image"`
	BaseModel
}
