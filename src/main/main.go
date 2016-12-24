package main

import (
	"fmt"
	"net/http"
	"html/template"
)

func index(response http.ResponseWriter, request *http.Request)  {
	templ, err := template.ParseFiles("templates/base.html")
	if err != nil {
		fmt.Fprintf(response, err.Error())
	}
	templ.ExecuteTemplate(response, "index", nil)
}

func main() {

	http.HandleFunc("/", index)
	http.ListenAndServe(":3000", nil)
}
