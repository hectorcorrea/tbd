package common

import (
	"net/http"
)

type Session struct {
	Resp      http.ResponseWriter
	Req       *http.Request
	Cookie    *http.Cookie
	LoginName string
	SessionId string
	ExpiresOn string
}

func NewSession(resp http.ResponseWriter, req *http.Request) Session {
	sessionId := ""
	login := ""
	expiresOn := ""
	cookie, _ := req.Cookie("sessionId")
	s := Session{
		Resp:      resp,
		Req:       req,
		Cookie:    cookie,
		LoginName: login,
		SessionId: sessionId,
		ExpiresOn: expiresOn,
	}
	return s
}
