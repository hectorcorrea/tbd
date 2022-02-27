package web

import (
	"net/http"
)

type session struct {
	resp      http.ResponseWriter
	req       *http.Request
	cookie    *http.Cookie
	loginName string
	sessionId string
	expiresOn string
}

func newSession(resp http.ResponseWriter, req *http.Request) session {
	sessionId := ""
	login := ""
	expiresOn := ""
	cookie, _ := req.Cookie("sessionId")
	s := session{
		resp:      resp,
		req:       req,
		cookie:    cookie,
		loginName: login,
		sessionId: sessionId,
		expiresOn: expiresOn,
	}
	return s
}
