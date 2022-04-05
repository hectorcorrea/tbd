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
	router.Add("POST", "/:id/post", post)
	router.Add("POST", "/:id/draft", draft)
	router.Add("POST", "/:id/edit", edit)
	router.Add("POST", "/:id/save", save)
	router.Add("POST", "/new", new)
	router.Add("GET", "/", viewAll)
	router.Add("GET", "/:slug/:id", viewOne)
}

func StartWebServer(address string, dataFolder string) {
	log.Printf("Listening for requests at %s\n", "http://"+address)
	db = textodb.InitTextoDb(dataFolder)
	http.HandleFunc("/", router.Dispatcher)

	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal("Failed to start the web server: ", err)
	}
}

func viewAll(resp http.ResponseWriter, req *http.Request, values map[string]string) {
	vm := db.All()
	renderTemplate(resp, req, "views/all.html", vm)
}

func viewOne(resp http.ResponseWriter, req *http.Request, values map[string]string) {
	id := values["id"]
	entry, err := db.FindById(id)
	if err != nil {
		log.Printf("Not found: %s. Error: %s", id, err)
		renderTemplate(resp, req, "views/error.html", entry)
		return
	}
	renderTemplate(resp, req, "views/one.html", entry)
}

func edit(resp http.ResponseWriter, req *http.Request, values map[string]string) {
	id := values["id"]
	entry, err := db.FindById(id)
	if err != nil {
		log.Printf("Not found: %s. Error: %s", id, err)
		renderTemplate(resp, req, "views/error.html", entry)
		return
	}
	renderTemplate(resp, req, "views/edit.html", entry)
}

func new(resp http.ResponseWriter, req *http.Request, values map[string]string) {
	entry, err := db.NewEntry()
	if err != nil {
		log.Printf("Error creating new document: %s", err)
		http.Error(resp, "Error processing request", http.StatusInternalServerError)
		return
	}

	renderAfterSave(resp, req, entry, "Created", nil)
}

func save(resp http.ResponseWriter, req *http.Request, values map[string]string) {
	id := values["id"]
	entry, err := db.FindById(id)
	if err != nil {
		log.Printf("Error fetching document to save: %s", err)
		http.Error(resp, "Error processing request", http.StatusInternalServerError)
		return
	}

	entry.Title = req.FormValue("title")
	entry.Summary = req.FormValue("summary")
	entry.SetContent(req.FormValue("content"))
	entry, err = db.UpdateEntry(entry)
	renderAfterSave(resp, req, entry, "Saved", err)
}

func post(resp http.ResponseWriter, req *http.Request, values map[string]string) {
	id := values["id"]
	entry, err := db.FindById(id)
	if err != nil {
		log.Printf("Error fetching document to save: %s", err)
		http.Error(resp, "Error processing request", http.StatusInternalServerError)
		return
	}

	entry.MarkAsPosted()
	entry, err = db.UpdateEntry(entry)
	renderAfterSave(resp, req, entry, "Posted", err)
}

func draft(resp http.ResponseWriter, req *http.Request, values map[string]string) {
	id := values["id"]
	entry, err := db.FindById(id)
	if err != nil {
		log.Printf("Error fetching document to save: %s", err)
		http.Error(resp, "Error processing request", http.StatusInternalServerError)
		return
	}

	entry.MarkAsDraft()
	entry, err = db.UpdateEntry(entry)
	renderAfterSave(resp, req, entry, "Marked as Draft", err)
}

func renderAfterSave(resp http.ResponseWriter, req *http.Request, entry textodb.TextoEntry, action string, err error) {
	if err != nil {
		log.Printf("Error saving document: %s", err)
		http.Error(resp, "Error processing request", http.StatusInternalServerError)
		return
	}

	qs := req.URL.Query()
	if len(qs["redirect"]) > 0 {
		// Redirect to view page
		url := fmt.Sprintf("/%s/%s", entry.Slug, entry.Id)
		log.Printf("%s %s, redirecting to %s", action, entry.Id, url)
		http.Redirect(resp, req, url, 301)
		return
	}

	// Return JSON payload
	log.Printf("%s %s", action, entry.Id)
	payload := fmt.Sprintf(`{ "id": "%s", "slug": "%s"}`, entry.Id, entry.Slug)
	resp.Header().Add("Content-Type", "text/json")
	fmt.Fprint(resp, payload)
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
