package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/codered-by-ec-council/Hands-on-Network-Programming-with-Go/pkg/models"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

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

func (app *application) handleVendorUpdate(w http.ResponseWriter, r *http.Request) {
	log.Println("hit vendor update handler")
	if err := r.ParseForm(); err != nil {
		http.Error(w, errors.New("couldn't parse form").Error(), http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, errors.New("id is not an int").Error(), http.StatusBadRequest)
		return
	}
	vendor := models.Vendor{
		ID:   id,
		Name: r.FormValue("vendor"),
	}
	if err := app.inventoryService.VendorRepo.Update(vendor); err != nil {
		http.Error(w, errors.New("couldn't update vendor in the DB").Error(), http.StatusBadRequest)
		return
	}
	vendors, err := app.inventoryService.VendorRepo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	app.render(w, "vendors.page.tmpl", vendors)
}

func (app *application) handleVendor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, errors.New("id is not an int").Error(), http.StatusBadRequest)
		return
	}
	vendor, err := app.inventoryService.VendorRepo.GetByID(id)
	if err != nil {
		http.Error(w, errors.New("vendor is not in DB").Error(), http.StatusBadRequest)
		return
	}
	app.render(w, "vendors_update.page.tmpl", vendor)
}

func (app *application) handleVendors(w http.ResponseWriter, r *http.Request) {
	vendors, err := app.inventoryService.VendorRepo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	app.render(w, "vendors.page.tmpl", vendors)
}

func (app *application) handleVendorDelete(w http.ResponseWriter, r *http.Request) {
	log.Println("hit vendor delete handler")
	if err := r.ParseForm(); err != nil {
		http.Error(w, errors.New("couldn't parse form").Error(), http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, errors.New("id is not an int").Error(), http.StatusBadRequest)
		return
	}
	if err := app.inventoryService.VendorRepo.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

func (app *application) handleModels(w http.ResponseWriter, r *http.Request) {
	models, err := app.inventoryService.ModelRepo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	app.render(w, "models.page.tmpl", models)
}

func (app *application) handleModel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, errors.New("id is not an int").Error(), http.StatusBadRequest)
		return
	}
	model, err := app.inventoryService.ModelRepo.GetByID(id)
	if err != nil {
		http.Error(w, errors.New("vendor is not in DB").Error(), http.StatusBadRequest)
		return
	}
	vendors, err := app.inventoryService.VendorRepo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	composite := struct {
		Vendors []models.Vendor
		Model   models.Model
	}{
		Vendors: vendors,
		Model:   model,
	}
	app.render(w, "models_update.page.tmpl", composite)
}

func (app *application) handleModelUpdate(w http.ResponseWriter, r *http.Request) {
	log.Println("hit vendor update handler")
	if err := r.ParseForm(); err != nil {
		http.Error(w, errors.New("couldn't parse form").Error(), http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, errors.New("id is not an int").Error(), http.StatusBadRequest)
		return
	}
	vendorID, err := strconv.Atoi(r.FormValue("vendor_id"))
	if err != nil {
		http.Error(w, errors.New("id is not an int").Error(), http.StatusBadRequest)
		return
	}
	model := models.Model{
		ID:     id,
		Name:   r.FormValue("model"),
		Vendor: models.Vendor{ID: vendorID},
	}
	if err := app.inventoryService.ModelRepo.Update(model); err != nil {
		http.Error(w, errors.New("couldn't update vendor in the DB").Error(), http.StatusBadRequest)
		return
	}
	models, err := app.inventoryService.ModelRepo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	app.render(w, "models.page.tmpl", models)
}

func (app *application) handleModelDelete(w http.ResponseWriter, r *http.Request) {
	log.Println("hit model delete handler")
	if err := r.ParseForm(); err != nil {
		http.Error(w, errors.New("couldn't parse form").Error(), http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, errors.New("id is not an int").Error(), http.StatusBadRequest)
		return
	}
	if err := app.inventoryService.ModelRepo.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

func (app *application) handleDevices(w http.ResponseWriter, r *http.Request) {
	devices, err := app.inventoryService.DeviceRepo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	app.render(w, "devices.page.tmpl", devices)
}

func (app *application) handleDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, errors.New("id is not an int").Error(), http.StatusBadRequest)
		return
	}
	device, err := app.inventoryService.DeviceRepo.GetByID(id)
	if err != nil {
		http.Error(w, errors.New("device is not in DB").Error(), http.StatusBadRequest)
		return
	}
	devModels, err := app.inventoryService.ModelRepo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	composite := struct {
		Models []models.Model
		Device models.Device
	}{
		Models: devModels,
		Device: device,
	}
	app.render(w, "devices_update.page.tmpl", composite)
}

func (app *application) handleDeviceUpdate(w http.ResponseWriter, r *http.Request) {
	log.Println("hit device update handler")
	if err := r.ParseForm(); err != nil {
		http.Error(w, errors.New("couldn't parse form").Error(), http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, errors.New("id is not an int").Error(), http.StatusBadRequest)
		return
	}
	modelID, err := strconv.Atoi(r.FormValue("model_id"))
	if err != nil {
		http.Error(w, errors.New("id is not an int").Error(), http.StatusBadRequest)
		return
	}
	device := models.Device{
		ID:       id,
		Hostname: r.FormValue("hostname"),
		IPv4:     r.FormValue("ip"),
		Model:    models.Model{ID: modelID},
	}
	if err := app.inventoryService.DeviceRepo.Update(device); err != nil {
		http.Error(w, errors.New("couldn't update vendor in the DB").Error(), http.StatusBadRequest)
		return
	}
	devices, err := app.inventoryService.DeviceRepo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	app.render(w, "devices.page.tmpl", devices)
}

func (app *application) handleDeviceDelete(w http.ResponseWriter, r *http.Request) {
	log.Println("hit device delete handler")
	if err := r.ParseForm(); err != nil {
		http.Error(w, errors.New("couldn't parse form").Error(), http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, errors.New("id is not an int").Error(), http.StatusBadRequest)
		return
	}
	if err := app.inventoryService.DeviceRepo.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	devices, err := app.inventoryService.DeviceRepo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	app.render(w, "devices.page.tmpl", devices)
}
