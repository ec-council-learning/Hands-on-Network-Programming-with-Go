package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/pkg/errors"
)

//go:embed templates/*
var content embed.FS

var templateCache = map[string]*template.Template{}

func main() {
	websocket := flag.String("websocket", "localhost:8080", "socket on which to listen for incoming connections")
	tc, err := newTemplateCache("templates")
	if err != nil {
		log.Println(err)
	}
	templateCache = tc
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/vendors", handleVendors)
	log.Println("starting web server on", *websocket)
	http.ListenAndServe(*websocket, nil)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	ts, ok := templateCache["home.page.tmpl"]
	if !ok {
		http.Error(w, fmt.Errorf("The template %s does not exist", "home").Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(ts.Name())
	fmt.Println(ts.DefinedTemplates())
	if err := ts.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleVendors(w http.ResponseWriter, r *http.Request) {
	ts, ok := templateCache["vendors.page.tmpl"]
	if !ok {
		http.Error(w, fmt.Errorf("The template %s does not exist", "vendors").Error(), http.StatusInternalServerError)
		return
	}
	if err := ts.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var functions = template.FuncMap{}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := content.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrap(err, "ReadDir failed")
	}
	for _, page := range pages {
		name := filepath.Base(page.Name())
		ts, err := template.New(name).Funcs(functions).ParseFS(content, fmt.Sprintf("templates/%v", name))
		if err != nil {
			return nil, errors.Wrap(err, "ParseFS failed")
		}
		ts, err = ts.ParseFS(content, "templates/*.layout.tmpl")
		if err != nil {
			return nil, errors.Wrap(err, "ParseFS failed")
		}
		cache[name] = ts
	}
	return cache, nil
}

func render(w http.ResponseWriter, r *http.Request, name string) {
	ts, ok := templateCache[name]
	if !ok {
		http.Error(w, fmt.Errorf("The template %s does not exist", name).Error(), http.StatusInternalServerError)
		return
	}
	if err := ts.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
