package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/hectorcorrea/texto/common"
	"github.com/hectorcorrea/texto/textdb"
)

var router common.Router
var db textdb.TextDb

func init() {
	router.Add("POST", "/doc/new", docNew)
	router.Add("GET", "/doc", docAll)
	router.Add("GET", "/doc/:slug", docOne)
}

func docAll(s common.Session, values map[string]string) {
	vm := db.ListAll()
	renderTemplate(s, "views/all.html", vm)
}

func docOne(s common.Session, values map[string]string) {
	slug := values["slug"]
	found, entry := db.FindBySlug(slug)
	if !found {
		log.Printf("Not found: %s", slug)
		renderTemplate(s, "views/error.html", entry)
		return
	}
	renderTemplate(s, "views/one.html", entry)
}

func docNew(s common.Session, values map[string]string) {
	entry, err := db.NewEntry()
	if err != nil {
		log.Printf("Error creating new document: %s", err)
		http.Error(s.Resp, "Error processing request", http.StatusInternalServerError)
		return
	}

	qs := s.Req.URL.Query()
	if len(qs["redirect"]) > 0 {
		url := fmt.Sprintf("/doc/%s", entry.Id)
		log.Printf("Created %s, redirecting to %s", entry.Id, url)
		http.Redirect(s.Resp, s.Req, url, 301)
		return
	}

	log.Printf("Created %s", entry.Id)
	payload := "{ \"path\":\"" + entry.Id + "\" }"
	s.Resp.Header().Add("Content-Type", "text/json")
	fmt.Fprint(s.Resp, payload)
}

func dispatcher(resp http.ResponseWriter, req *http.Request) {
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
	http.HandleFunc("/doc/", dispatcher)

	err := http.ListenAndServe(settings.Address, nil)
	if err != nil {
		log.Fatal("Failed to start the web server: ", err)
	}
}

func renderTemplate(s common.Session, viewName string, viewModel interface{}) {
	t, err := loadTemplate(s, viewName)
	if err != nil {
		log.Printf("Error loading: %s, %s ", viewName, err)
	} else {
		err = t.Execute(s.Resp, viewModel)
		if err != nil {
			log.Printf("Error rendering: %s, %s ", viewName, err)
		}
	}
}

func loadTemplate(s common.Session, viewName string) (*template.Template, error) {
	t, err := template.New("layout").ParseFiles("views/layout.html", viewName)
	if err != nil {
		log.Printf("Error loading template %s (%s)", viewName, s.Req.URL.Path)
		return nil, err
	} else {
		log.Printf("Loaded template %s (%s)", viewName, s.Req.URL.Path)
		return t, nil
	}
}
