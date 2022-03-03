package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hectorcorrea/texto/common"
	"github.com/hectorcorrea/texto/textdb"
)

var router common.Router
var db textdb.TextDb

func init() {
	router.Add("POST", "/record/new", recordNew)
	router.Add("GET", "/record", recordAll)
	router.Add("GET", "/record/:id", recordOne)
}

func recordAll(s common.Session, values map[string]string) {
	fmt.Fprint(s.Resp, "record all")
}

func recordOne(s common.Session, values map[string]string) {
	fmt.Fprint(s.Resp, "record id")
}

func recordNew(s common.Session, values map[string]string) {
	err := db.CreateNewEntry()
	if err != nil {
		fmt.Fprint(s.Resp, err)
	}
	fmt.Fprint(s.Resp, "created new")
}

func recordPages(resp http.ResponseWriter, req *http.Request) {
	session := common.NewSession(resp, req)
	found, route := router.FindRoute(req.Method, req.URL.Path)
	if found {
		values := route.UrlValues(req.URL.Path)
		route.Handler(session, values)
	} else {
		log.Printf("not found")
	}
}

func StartWebServer(settings Settings) {
	log.Printf("Listening for requests at %s\n", "http://"+settings.Address)
	db = textdb.InitTextDb(settings.DataFolder)
	http.HandleFunc("/record/", recordPages)

	err := http.ListenAndServe(settings.Address, nil)
	if err != nil {
		log.Fatal("Failed to start the web server: ", err)
	}
}
