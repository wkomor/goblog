package main

import (
	"fmt"
	"net/http"
	"html/template"
	"models"
)

var posts map[string] *models.Post

func index_handler(response http.ResponseWriter, request *http.Request)  {
	templ, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(response, err.Error())
	}
	fmt.Println(posts)
	templ.ExecuteTemplate(response, "index", posts)
}

func write_handler(response http.ResponseWriter, request *http.Request)  {
	templ, err := template.ParseFiles("templates/write.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(response, err.Error())
	}
	templ.ExecuteTemplate(response, "write", nil)
}

func save_post_handler(response http.ResponseWriter, request *http.Request)  {
	id := GenerateId()
	title := request.FormValue("title")
	content := request.FormValue("content")
	post := models.New_Post(id, title, content)
	posts[post.Id] = post
	http.Redirect(response, request, "/", 302)
}


func main() {

	posts = make(map[string]*models.Post, 0)

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))
	http.HandleFunc("/", index_handler)
	http.HandleFunc("/write", write_handler)
	http.HandleFunc("/save_post", save_post_handler)
	http.ListenAndServe(":3000", nil)
}
