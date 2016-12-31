package models

type Post struct {
	Id          string
	Title       string
	ContentHTML string
	ContentMD   string
}

func New_Post(id, title, content_html, content_md string) *Post {
	return &Post{id, title, content_html, content_md}
}
