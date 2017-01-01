package main

import (
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/martini-contrib/render"
	"html/template"
	"models"
	"net/http"
	"labix.org/v2/mgo"
	"documents"
	"utils"
	"session"
	"time"
)
const (
	COOKIE_NAME = "sessionId"
)
var postsCollection *mgo.Collection
var inMemorySession *session.Session

func index_handler(rnd render.Render, request *http.Request) {
	cookie, _ := request.Cookie(COOKIE_NAME)
	if cookie != nil {
		fmt.Println(inMemorySession.Get(cookie.Value))
	}
	postDocuments := []documents.PostDocument{}
	postsCollection.Find(nil).All(&postDocuments)
	posts := []models.Post{}
	for _, doc := range postDocuments{
		post := models.Post{doc.Id, doc.Title, doc.ContentHTML, doc.ContentMD}
		posts = append(posts, post)
	}

	rnd.HTML(200, "index", posts)
}

func write_handler(rnd render.Render) {
	rnd.HTML(200, "write", nil)
}

func edit_handler(rnd render.Render, params martini.Params) {
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

func save_post_handler(rnd render.Render, request *http.Request) {
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

func unescape(x string) interface{} {
	return template.HTML(x)
}

func delete_handler(rnd render.Render, request *http.Request, params martini.Params) {
	id := params["id"]
	if id == "" {
		rnd.Redirect("/")
		return
	}

	postsCollection.RemoveId(id)

	rnd.Redirect("/")
}

func getHTML(rnd render.Render, request *http.Request) {
	md := request.FormValue("md")
	html := utils.MarcdownToHTML(md)
	rnd.JSON(200, map[string]interface{}{"html": string(html)})
}

func getLoginHandler(rnd render.Render){
	rnd.HTML(200, "login", nil)
}

func postLoginHandler(rnd render.Render, response http.ResponseWriter, request *http.Request){
	username := request.FormValue("username")
	password := request.FormValue("password")

	sessionId := inMemorySession.Init(username)
	fmt.Println(password)

	cookie := &http.Cookie{
		Name: COOKIE_NAME,
		Value:sessionId,
		Expires: time.Now().Add(5 * time.Minute),
	}

	http.SetCookie(response, cookie)

	rnd.Redirect("/", 302)
}

func main() {

	sessionDB, err := mgo.Dial("localhost")

	if err!=nil{
		panic(err)
	}

	inMemorySession = session.NewSession()

	postsCollection = sessionDB.DB("blog").C("posts")

	m := martini.Classic()

	unescapeFuncMap := template.FuncMap{"unescape": unescape}

	m.Use(render.Renderer(render.Options{
		Directory:  "templates",                         // Specify what path to load the templates from.
		Layout:     "layout",                            // Specify a layout template. Layouts can call {{ yield }} to render the current template.
		Extensions: []string{".tmpl", ".html"},          // Specify extensions to load for templates.
		Funcs:      []template.FuncMap{unescapeFuncMap}, // Specify helper function maps for templates to access.
		Charset:    "UTF-8",                             // Sets encoding for json and html content-types. Default is "UTF-8".
		IndentJSON: true,                                // Output human readable JSON
	}))

	staticOptions := martini.StaticOptions{Prefix: "assets"}
	m.Use(martini.Static("assets", staticOptions))
	m.Get("/", index_handler)
	m.Get("/write", write_handler)
	m.Get("/login", getLoginHandler)
	m.Post("/login", postLoginHandler)
	m.Get("/edit/:id", edit_handler)
	m.Get("/delete/:id", delete_handler)
	m.Post("/save_post", save_post_handler)
	m.Post("/gethtml", getHTML)

	m.Get("/test", func() string {
		return "test"
	})

	m.Run()
}
