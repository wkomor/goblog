package session

import (
	"utils"
	"github.com/codegangsta/martini"
	"net/http"
	"time"
)

const (
	COOKIE_NAME = "sessionId"
)

type Session struct {
	id       string
	Username string
	IsAutorized bool
}

type SessionStore struct {
	data map[string]*Session
}

func NewSessionStore() *SessionStore {
	s := new(SessionStore)
	s.data = make(map[string]*Session)
	return s
}

func (store *SessionStore) Get(sessionId string) *Session  {
	session := store.data[sessionId]
	if session == nil {
		return &Session{id: sessionId}
	}
	return session
}

func (store *SessionStore) Set(session *Session) {
	store.data[session.id] = session
}

func ensureCookie(request *http.Request, response http.ResponseWriter) string {
	cookie, _ := request.Cookie(COOKIE_NAME)
	if cookie != nil{
		return cookie.Value
	}
	sessionId := utils.GenerateId()
	cookie = &http.Cookie{Name:COOKIE_NAME,
		              Value:sessionId,
			      Expires:time.Now().Add(5*time.Minute),
	}
	http.SetCookie(response, cookie)
	return sessionId

}

var sessionStore = NewSessionStore()

func Middleware(ctx martini.Context, request *http.Request, response http.ResponseWriter)  {
	sessionId := ensureCookie(request, response)
	session := sessionStore.Get(sessionId)
	ctx.Map(session)
	ctx.Next()
	sessionStore.Set(session)
}