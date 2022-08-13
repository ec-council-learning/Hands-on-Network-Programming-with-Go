package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/codered-by-ec-council/Hands-on-Network-Programming-with-Go/pkg/inventory"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

//go:embed templates/*
var content embed.FS

type application struct {
	templateCache    map[string]*template.Template
	inventoryService *inventory.Service
}

func main() {
	websocket := flag.String("websocket", "localhost:8080", "socket on which to listen for incoming connections")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	dbpool, err := pgxpool.Connect(ctx, os.Getenv("DSN"))
	if err != nil {
		log.Println(err)
	}
	defer dbpool.Close()
	if err := dbpool.Ping(ctx); err != nil {
		log.Println("db conn failed:", err)
	}
	fmt.Println("db connected successfully")
	tc, err := newTemplateCache("templates")
	if err != nil {
		log.Println(err)
	}
	app := application{
		templateCache:    tc,
		inventoryService: inventory.NewService(dbpool),
	}
	http.HandleFunc("/", app.handleHome)
	http.HandleFunc("/vendors", app.handleVendors)
	log.Println("starting web server on", *websocket)
	http.ListenAndServe(*websocket, nil)
}

func (app *application) handleHome(w http.ResponseWriter, r *http.Request) {
	app.render(w, "home.page.tmpl")
}

func (app *application) handleVendors(w http.ResponseWriter, r *http.Request) {
	app.render(w, "vendors.page.tmpl")
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

func (app *application) render(w http.ResponseWriter, name string) {
	ts, ok := app.templateCache[name]
	if !ok {
		http.Error(w, fmt.Errorf("The template %s does not exist", name).Error(), http.StatusInternalServerError)
		return
	}
	if err := ts.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
