package routes

import (
	"github.com/martini-contrib/render"
	"documents"
	"models"
	"labix.org/v2/mgo"
	"session"
)

func IndexHandler(rnd render.Render, db *mgo.Database, s *session.Session) {
	postsCollection := db.C("posts")

	postDocuments := []documents.PostDocument{}
	postsCollection.Find(nil).All(&postDocuments)
	posts := []models.Post{}
	for _, doc := range postDocuments{
		post := models.Post{doc.Id, doc.Title, doc.ContentHTML, doc.ContentMD}
		posts = append(posts, post)
	}

	model := models.PostListModel{}
	model.IsAutorized = s.IsAutorized
	model.Posts = posts

	rnd.HTML(200, "index", model)
}
