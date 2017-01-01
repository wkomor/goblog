package session

import "utils"

type sessionData struct {
	Username string
}

type Session struct {
	data map[string]*sessionData
}

func NewSession() *Session {
	s := new(Session)
	s.data = make(map[string]*sessionData)
	return s
}

func (self *Session) Init(username string) string {
	id:= utils.GenerateId()
	data := &sessionData{Username: username}
	self.data[id] = data
	return id

}

func (self *Session) Get(sessionId string) string {
	data := self.data[sessionId]
	if data == nil{
		return ""
	}
	return data.Username
}