package main

import (
	"fmt"
	"net/http"
	"html/template"
	"github.com/codegangsta/martini"
	"github.com/martini-contrib/render"
	"models"
)

var posts map[string] *models.Post
var counter int

func index_handler(rnd render.Render, response http.ResponseWriter, request *http.Request)  {
	templ, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(response, err.Error())
	}
	fmt.Println(counter)
	templ.ExecuteTemplate(response, "index", posts)
}

func write_handler(rnd render.Render, response http.ResponseWriter, request *http.Request)  {
	templ, err := template.ParseFiles("templates/write.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(response, err.Error())
	}
	templ.ExecuteTemplate(response, "write", nil)
}

func edit_handler(rnd render.Render, response http.ResponseWriter, request *http.Request)  {
	templ, err := template.ParseFiles("templates/write.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(response, err.Error())
	}
	id := request.FormValue("id")
	post, found := posts[id]
	if !found {
		http.NotFound(response, request)
	}
	templ.ExecuteTemplate(response, "write", post)
}

func save_post_handler(response http.ResponseWriter, request *http.Request)  {
	id := request.FormValue("id")
	title := request.FormValue("title")
	content := request.FormValue("content")

	var post *models.Post
	if id!=""{
		post = posts[id]
		post.Title = title
		post.Content = content
	} else {
		id := GenerateId()
		post := models.New_Post(id, title, content)
		posts[post.Id] = post
	}

	http.Redirect(response, request, "/", 302)
}


func delete_handler(response http.ResponseWriter, request *http.Request)  {
	id := request.FormValue("id")
	if id == "" {
		http.NotFound(response, request)
	}

	delete(posts, id)

	http.Redirect(response, request, "/", 302)
}

func main() {

	posts = make(map[string]*models.Post, 0)

	m := martini.Classic()

	counter = 0

	m.Use(render.Renderer(render.Options{
	  Directory: "templates", // Specify what path to load the templates from.
	  Layout: "layout", // Specify a layout template. Layouts can call {{ yield }} to render the current template.
	  Extensions: []string{".tmpl", ".html"}, // Specify extensions to load for templates.
	  //Funcs: []template.FuncMap{AppHelpers}, // Specify helper function maps for templates to access.
	  Charset: "UTF-8", // Sets encoding for json and html content-types. Default is "UTF-8".
	  IndentJSON: true, // Output human readable JSON
	  HTMLContentType: "application/xhtml+xml", // Output XHTML content type instead of default "text/html"
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
	m.Get("/edit", edit_handler)
	m.Get("/delete", delete_handler)
	m.Post("/save_post", save_post_handler)

	m.Get("/test", func() string {
		return "test"
	})

	m.Run()
}
