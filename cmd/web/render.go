package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/pkg/errors"
)

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

func (app *application) render(w http.ResponseWriter, name string, data interface{}) {
	ts, ok := app.templateCache[name]
	if !ok {
		http.Error(w, fmt.Errorf("The template %s does not exist", name).Error(), http.StatusInternalServerError)
		return
	}
	if err := ts.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
