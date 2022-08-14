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
	"strconv"
	"text/template"
	"time"

	"github.com/codered-by-ec-council/Hands-on-Network-Programming-with-Go/pkg/inventory"
	"github.com/codered-by-ec-council/Hands-on-Network-Programming-with-Go/pkg/models"
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
	http.HandleFunc("/vendors/new", app.handleVendorsNew)
	http.HandleFunc("/vendors/add", app.handleVendorsAdd)
	http.HandleFunc("/models", app.handleModels)
	http.HandleFunc("/models/new", app.handleModelsNew)
	http.HandleFunc("/models/add", app.handleModelsAdd)
	http.HandleFunc("/devices/new", app.handleDeviceNew)
	http.HandleFunc("/devices/add", app.handleDeviceAdd)
	http.HandleFunc("/devices", app.handleDevices)
	log.Println("starting web server on", *websocket)
	http.ListenAndServe(*websocket, nil)
}

func (app *application) handleHome(w http.ResponseWriter, r *http.Request) {
	app.render(w, "home.page.tmpl", nil)
}

func (app *application) handleVendorsNew(w http.ResponseWriter, r *http.Request) {
	app.render(w, "vendors_new.page.tmpl", nil)
}

func (app *application) handleVendorsAdd(w http.ResponseWriter, r *http.Request) {
	log.Println("hit vendor add handler")
	if r.Method != http.MethodPost {
		http.Error(w, errors.New("method should be post").Error(), http.StatusBadRequest)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, errors.New("couldn't parse form").Error(), http.StatusBadRequest)
		return
	}
	if err := app.inventoryService.VendorRepo.New(r.FormValue("vendor")); err != nil {
		http.Error(w, errors.New("couldn't add vendor to DB").Error(), http.StatusBadRequest)
		return
	}
	vendors, err := app.inventoryService.VendorRepo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	app.render(w, "vendors.page.tmpl", vendors)
}

func (app *application) handleModelsNew(w http.ResponseWriter, r *http.Request) {
	vendors, err := app.inventoryService.VendorRepo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	app.render(w, "models_new.page.tmpl", vendors)
}

func (app *application) handleModelsAdd(w http.ResponseWriter, r *http.Request) {
	log.Println("hit model add handler")
	if r.Method != http.MethodPost {
		http.Error(w, errors.New("method should be post").Error(), http.StatusBadRequest)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, errors.New("couldn't parse form").Error(), http.StatusBadRequest)
		return
	}
	idStr := r.FormValue("vendor_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, errors.New("id is not an int").Error(), http.StatusBadRequest)
		return
	}
	newModel := models.Model{
		Name:   r.FormValue("model"),
		Vendor: models.Vendor{ID: id},
	}
	if err := app.inventoryService.ModelRepo.New(newModel); err != nil {
		http.Error(w, errors.New("couldn't add model to DB").Error(), http.StatusBadRequest)
		return
	}
	models, err := app.inventoryService.ModelRepo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	app.render(w, "models.page.tmpl", models)
}

func (app *application) handleDeviceNew(w http.ResponseWriter, r *http.Request) {
	models, err := app.inventoryService.ModelRepo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	app.render(w, "devices_new.page.tmpl", models)
}

func (app *application) handleDeviceAdd(w http.ResponseWriter, r *http.Request) {
	log.Println("hit device add handler")
	if r.Method != http.MethodPost {
		http.Error(w, errors.New("method should be post").Error(), http.StatusBadRequest)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, errors.New("couldn't parse form").Error(), http.StatusBadRequest)
		return
	}
	idStr := r.FormValue("model_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, errors.New("id is not an int").Error(), http.StatusBadRequest)
		return
	}
	newDevice := models.Device{
		Hostname: r.FormValue("hostname"),
		IPv4:     r.FormValue("ip"),
		Model:    models.Model{ID: id},
	}
	if err := app.inventoryService.DeviceRepo.New(newDevice); err != nil {
		http.Error(w, errors.New("couldn't add model to DB").Error(), http.StatusBadRequest)
		return
	}
	devices, err := app.inventoryService.DeviceRepo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	app.render(w, "devices.page.tmpl", devices)
}

func (app *application) handleVendors(w http.ResponseWriter, r *http.Request) {
	vendors, err := app.inventoryService.VendorRepo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	app.render(w, "vendors.page.tmpl", vendors)
}

func (app *application) handleModels(w http.ResponseWriter, r *http.Request) {
	models, err := app.inventoryService.ModelRepo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	app.render(w, "models.page.tmpl", models)
}

func (app *application) handleDevices(w http.ResponseWriter, r *http.Request) {
	devices, err := app.inventoryService.DeviceRepo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	app.render(w, "devices.page.tmpl", devices)
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
