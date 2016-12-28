package models

type Post struct {
	Id      string
	Title   string
	Content string
}

func New_Post(id, title, content string) *Post {
	return &Post{id, title, content}
}
