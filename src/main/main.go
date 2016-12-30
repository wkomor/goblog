package main

import (
	"fmt"
	"net/http"
	"github.com/codegangsta/martini"
	"github.com/martini-contrib/render"
	"github.com/russross/blackfriday"
	"models"
	"html/template"
)

var posts map[string] *models.Post
var counter int

func index_handler(rnd render.Render)  {
	fmt.Println(counter)
	rnd.HTML(200, "index", posts)
}

func write_handler(rnd render.Render)  {
	rnd.HTML(200, "write", nil)
}

func edit_handler(rnd render.Render, request *http.Request, params martini.Params)  {
	id := params["id"]
	post, found := posts[id]
	if !found {
		rnd.Redirect("/")
		return
	}
	rnd.HTML(200, "write", post)
}

func save_post_handler(rnd render.Render, request *http.Request)  {
	id := request.FormValue("id")
	title := request.FormValue("title")
	content_markdown := request.FormValue("content")

	content_html := string(blackfriday.MarkdownBasic([]byte(content_markdown)))

	var post *models.Post
	if id!=""{
		post = posts[id]
		post.Title = title
		post.ContentHTML = content_html
		post.ContentMD = content_markdown
	} else {
		id := GenerateId()
		post := models.New_Post(id, title, content_html, content_markdown)
		posts[post.Id] = post
	}

	rnd.Redirect("/")
}

func unescape(x string) interface{} {
	return template.HTML(x)
}


func delete_handler(rnd render.Render, request *http.Request, params martini.Params)  {
	id := params["id"]
	if id == "" {
		rnd.Redirect("/")
		return
	}

	delete(posts, id)

	rnd.Redirect("/")
}

func getHTML(rnd render.Render, request *http.Request)  {
	md := request.FormValue("md")
	htmlBytes := blackfriday.MarkdownBasic([]byte(md))
	rnd.JSON(200, map[string] interface{} {"html": string(htmlBytes)})
}

func main() {

	posts = make(map[string]*models.Post, 0)

	m := martini.Classic()

	counter = 0

	unescapeFuncMap := template.FuncMap{"unescape": unescape}

	m.Use(render.Renderer(render.Options{
	  Directory: "templates", // Specify what path to load the templates from.
	  Layout: "layout", // Specify a layout template. Layouts can call {{ yield }} to render the current template.
	  Extensions: []string{".tmpl", ".html"}, // Specify extensions to load for templates.
	  Funcs: []template.FuncMap{unescapeFuncMap}, // Specify helper function maps for templates to access.
	  Charset: "UTF-8", // Sets encoding for json and html content-types. Default is "UTF-8".
	  IndentJSON: true, // Output human readable JSON
	}))

	staticOptions := martini.StaticOptions{Prefix:"assets"}
	m.Use(martini.Static("assets", staticOptions))
	m.Use(func(r *http.Request) {
		if r.URL.Path == "/write" {
			counter++
		}
	})
	m.Get("/", index_handler)
	m.Get("/write", write_handler)
	m.Get("/edit/:id", edit_handler)
	m.Get("/delete/:id", delete_handler)
	m.Post("/save_post", save_post_handler)
	m.Post("/gethtml", getHTML)

	m.Get("/test", func() string {
		return "test"
	})

	m.Run()
}
