package main

import (
	"github.com/codegangsta/martini"
	"github.com/martini-contrib/render"
	"html/template"
	"labix.org/v2/mgo"
	"session"
	"routes"
)

func unescape(x string) interface{} {
	return template.HTML(x)
}

func main() {

	sessionDB, err := mgo.Dial("localhost")

	if err!=nil{
		panic(err)
	}


	db := sessionDB.DB("blog")

	m := martini.Classic()

	m.Map(db)
	m.Use(session.Middleware)

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
	m.Get("/", routes.IndexHandler)
	m.Get("/write", routes.WriteHandler)
	m.Get("/login", routes.GetLoginHandler)
	m.Get("/logout", routes.GetLogoutHandler)
	m.Post("/login", routes.PostLoginHandler)
	m.Get("/edit/:id", routes.EditHandler)
	m.Get("/view/:id", routes.ViewHandler)
	m.Get("/delete/:id", routes.DeleteHandler)
	m.Post("/save_post", routes.SavePostsHandler)
	m.Post("/gethtml", routes.GetHTML)

	m.Get("/test", func() string {
		return "test"
	})

	m.Run()
}
