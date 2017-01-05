package routes

import (
	"github.com/martini-contrib/render"
	"net/http"
	"utils"
	"documents"
	"fmt"
	"models"
	"github.com/codegangsta/martini"
	"labix.org/v2/mgo"
)

func WriteHandler(rnd render.Render) {
	rnd.HTML(200, "write", nil)
}

func EditHandler(rnd render.Render, params martini.Params, db *mgo.Database) {
	postsCollection := db.C("posts")
	id := params["id"]
	postDocument := documents.PostDocument{}
	err := postsCollection.FindId(id).One(&postDocument)
	if err != nil {
		fmt.Println(err)
		rnd.Redirect("/")
		return
	}
	post := models.Post{postDocument.Id, postDocument.Title, postDocument.ContentHTML, postDocument.ContentMD}

	rnd.HTML(200, "write", post)
}

func SavePostsHandler(rnd render.Render, request *http.Request, db *mgo.Database) {
	postsCollection := db.C("posts")
	id := request.FormValue("id")
	title := request.FormValue("title")
	content_markdown := request.FormValue("content")

	content_html := utils.MarcdownToHTML(content_markdown)

	postDocument := documents.PostDocument{id, title, content_html, content_markdown}
	if id != "" {
		postsCollection.UpdateId(id, postDocument)

	} else {
		id := utils.GenerateId()
		postDocument.Id = id
		postsCollection.Insert(postDocument )
	}

	rnd.Redirect("/")
}


func DeleteHandler(rnd render.Render, params martini.Params, db *mgo.Database) {
	postsCollection := db.C("posts")
	id := params["id"]
	if id == "" {
		rnd.Redirect("/")
		return
	}

	postsCollection.RemoveId(id)

	rnd.Redirect("/")
}

func GetHTML(rnd render.Render, request *http.Request) {
	md := request.FormValue("md")
	html := utils.MarcdownToHTML(md)
	rnd.JSON(200, map[string]interface{}{"html": string(html)})
}
