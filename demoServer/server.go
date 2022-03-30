package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/hectorcorrea/textodb"
)

var router Router
var db textodb.TextoDb

func init() {
	router.Add("POST", "/doc/:id/edit", docEdit)
	router.Add("POST", "/doc/:id/save", docSave)
	router.Add("POST", "/doc/new", docNew)
	router.Add("GET", "/doc", docAll)
	router.Add("GET", "/doc/:slug", docOne)
}

func StartWebServer(address string, dataFolder string) {
	log.Printf("Listening for requests at %s\n", "http://"+address)
	db = textodb.InitTextDb(dataFolder)
	http.HandleFunc("/doc/", dispatcher)

	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal("Failed to start the web server: ", err)
	}
}

func docAll(resp http.ResponseWriter, req *http.Request, values map[string]string) {
	vm := db.All()
	renderTemplate(resp, req, "views/all.html", vm)
}

func docOne(resp http.ResponseWriter, req *http.Request, values map[string]string) {
	slug := values["slug"]
	entry, found := db.FindBySlug(slug)
	if !found {
		log.Printf("Not found: %s", slug)
		renderTemplate(resp, req, "views/error.html", entry)
		return
	}
	renderTemplate(resp, req, "views/one.html", entry)
}

func docEdit(resp http.ResponseWriter, req *http.Request, values map[string]string) {
	id := values["id"]
	log.Printf("id: %s.", id)

	entry, err := db.FindById(id)
	if err != nil {
		log.Printf("Not found id for edit: %s, %s", id, err)
		renderTemplate(resp, req, "views/error.html", entry)
		return
	}
	renderTemplate(resp, req, "views/edit.html", entry)
}

func docNew(resp http.ResponseWriter, req *http.Request, values map[string]string) {
	entry, err := db.NewEntry()
	if err != nil {
		log.Printf("Error creating new document: %s", err)
		http.Error(resp, "Error processing request", http.StatusInternalServerError)
		return
	}

	qs := req.URL.Query()
	if len(qs["redirect"]) > 0 {
		url := fmt.Sprintf("/doc/%s", entry.Slug)
		log.Printf("Created %s, redirecting to %s", entry.Id, url)
		http.Redirect(resp, req, url, 301)
		return
	}

	log.Printf("Created %s %s", entry.Id, entry.Slug)
	payload := "{ \"slug\":\"" + entry.Slug + "\" }"
	resp.Header().Add("Content-Type", "text/json")
	fmt.Fprint(resp, payload)
}

func docSave(resp http.ResponseWriter, req *http.Request, values map[string]string) {
	id := values["id"]
	entry, err := db.FindById(id)
	if err != nil {
		log.Printf("Error fetching document to save: %s", err)
		http.Error(resp, "Error processing request", http.StatusInternalServerError)
		return
	}

	entry.Title = req.FormValue("title")
	entry.Summary = req.FormValue("summary")
	entry.Content = req.FormValue("content")

	if req.FormValue("post") == "post" {
		entry.MarkAsPosted()
	} else if req.FormValue("draft") == "draft" {
		entry.MarkAsDraft()
	}

	entry, err = db.UpdateEntry(entry)
	if err != nil {
		log.Printf("Error saving document: %s", err)
		http.Error(resp, "Error processing request", http.StatusInternalServerError)
		return
	}

	qs := req.URL.Query()
	if len(qs["redirect"]) > 0 {
		url := fmt.Sprintf("/doc/%s", entry.Slug)
		log.Printf("Saved %s, redirecting to %s", entry.Id, url)
		http.Redirect(resp, req, url, 301)
		return
	}

	log.Printf("Saved %s", entry.Id)
	payload := "{ \"slug\":\"" + entry.Slug + "\" }"
	resp.Header().Add("Content-Type", "text/json")
	fmt.Fprint(resp, payload)
}

func dispatcher(resp http.ResponseWriter, req *http.Request) {
	found, route := router.FindRoute(req.Method, req.URL.Path)
	if found {
		values := route.UrlValues(req.URL.Path)
		route.Handler(resp, req, values)
	} else {
		log.Printf("not found")
	}
}

func renderTemplate(resp http.ResponseWriter, req *http.Request, viewName string, viewModel interface{}) {
	t, err := loadTemplate(resp, req, viewName)
	if err != nil {
		log.Printf("Error loading: %s, %s ", viewName, err)
	} else {
		err = t.Execute(resp, viewModel)
		if err != nil {
			log.Printf("Error rendering: %s, %s ", viewName, err)
		}
	}
}

func loadTemplate(resp http.ResponseWriter, req *http.Request, viewName string) (*template.Template, error) {
	t, err := template.New("layout").ParseFiles("views/layout.html", viewName)
	if err != nil {
		log.Printf("Error loading template %s (%s)", viewName, req.URL.Path)
		return nil, err
	} else {
		log.Printf("Loaded template %s (%s)", viewName, req.URL.Path)
		return t, nil
	}
}
