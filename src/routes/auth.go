package routes

import (
	"github.com/martini-contrib/render"
	"net/http"
	"fmt"
	"session"
)
func GetLoginHandler(rnd render.Render){
	rnd.HTML(200, "login", nil)
}

func GetLogoutHandler(rnd render.Render, s *session.Session){
	s.IsAutorized = false
	s.Username = ""
	rnd.Redirect("/", 302)
}

func PostLoginHandler(rnd render.Render, request *http.Request, s *session.Session){
	username := request.FormValue("username")
	password := request.FormValue("password")

	s.Username = username
	s.IsAutorized = true

	fmt.Println(password)


	rnd.Redirect("/", 302)
}


