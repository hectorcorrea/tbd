package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/hectorcorrea/tbd/textdb"
)

var router Router
var db textdb.TextDb

func init() {
	router.Add("POST", "/doc/:id/edit", docEdit)
	router.Add("POST", "/doc/:id/save", docSave)
	router.Add("POST", "/doc/new", docNew)
	router.Add("GET", "/doc", docAll)
	router.Add("GET", "/doc/:slug", docOne)
}

func docAll(s Session, values map[string]string) {
	vm := db.All()
	renderTemplate(s, "views/all.html", vm)
}

func docOne(s Session, values map[string]string) {
	slug := values["slug"]
	entry, found := db.FindBySlug(slug)
	if !found {
		log.Printf("Not found: %s", slug)
		renderTemplate(s, "views/error.html", entry)
		return
	}
	renderTemplate(s, "views/one.html", entry)
}

func docEdit(s Session, values map[string]string) {
	id := values["id"]
	log.Printf("id: %s.", id)

	entry, err := db.FindById(id)
	if err != nil {
		log.Printf("Not found id for edit: %s, %s", id, err)
		renderTemplate(s, "views/error.html", entry)
		return
	}
	renderTemplate(s, "views/edit.html", entry)
}

func docNew(s Session, values map[string]string) {
	entry, err := db.NewEntry()
	if err != nil {
		log.Printf("Error creating new document: %s", err)
		http.Error(s.Resp, "Error processing request", http.StatusInternalServerError)
		return
	}

	qs := s.Req.URL.Query()
	if len(qs["redirect"]) > 0 {
		url := fmt.Sprintf("/doc/%s", entry.Metadata.Slug)
		log.Printf("Created %s, redirecting to %s", entry.Id, url)
		http.Redirect(s.Resp, s.Req, url, 301)
		return
	}

	log.Printf("Created %s %s", entry.Id, entry.Metadata.Slug)
	payload := "{ \"slug\":\"" + entry.Metadata.Slug + "\" }"
	s.Resp.Header().Add("Content-Type", "text/json")
	fmt.Fprint(s.Resp, payload)
}

func docSave(s Session, values map[string]string) {
	id := values["id"]
	entry, err := db.FindById(id)
	if err != nil {
		log.Printf("Error fetching document to save: %s", err)
		http.Error(s.Resp, "Error processing request", http.StatusInternalServerError)
		return
	}

	entry.Metadata.Title = s.Req.FormValue("title")
	entry.Metadata.Slug = s.Req.FormValue("slug")
	entry.Metadata.Summary = s.Req.FormValue("summary")

	if s.Req.FormValue("post") == "post" {
		entry.MarkAsPosted()
	} else if s.Req.FormValue("draft") == "draft" {
		entry.MarkAsDraft()
	}

	entry, err = db.UpdateEntry(entry)
	if err != nil {
		log.Printf("Error saving document: %s", err)
		http.Error(s.Resp, "Error processing request", http.StatusInternalServerError)
		return
	}

	qs := s.Req.URL.Query()
	if len(qs["redirect"]) > 0 {
		url := fmt.Sprintf("/doc/%s", entry.Metadata.Slug)
		log.Printf("Saved %s, redirecting to %s", entry.Id, url)
		http.Redirect(s.Resp, s.Req, url, 301)
		return
	}

	log.Printf("Saved %s", entry.Id)
	payload := "{ \"slug\":\"" + entry.Metadata.Slug + "\" }"
	s.Resp.Header().Add("Content-Type", "text/json")
	fmt.Fprint(s.Resp, payload)
}

func dispatcher(resp http.ResponseWriter, req *http.Request) {
	session := NewSession(resp, req)
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

func renderTemplate(s Session, viewName string, viewModel interface{}) {
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

func loadTemplate(s Session, viewName string) (*template.Template, error) {
	t, err := template.New("layout").ParseFiles("views/layout.html", viewName)
	if err != nil {
		log.Printf("Error loading template %s (%s)", viewName, s.Req.URL.Path)
		return nil, err
	} else {
		log.Printf("Loaded template %s (%s)", viewName, s.Req.URL.Path)
		return t, nil
	}
}
