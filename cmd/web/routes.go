package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) routes() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/", app.handleHome).Methods(http.MethodGet)
	router.HandleFunc("/vendors", app.handleVendors).Methods(http.MethodGet)
	router.HandleFunc("/vendors/{id:[0-9]+}", app.handleVendor).Methods(http.MethodGet)
	router.HandleFunc("/vendors/{id:[0-9]+}", app.handleVendorUpdate).Methods(http.MethodPost)
	router.HandleFunc("/vendors/new", app.handleVendorsNew).Methods(http.MethodGet)
	router.HandleFunc("/vendors/add", app.handleVendorsAdd).Methods(http.MethodPost)
	router.HandleFunc("/models", app.handleModels).Methods(http.MethodGet)
	router.HandleFunc("/models/new", app.handleModelsNew).Methods(http.MethodGet)
	router.HandleFunc("/models/add", app.handleModelsAdd).Methods(http.MethodPost)
	router.HandleFunc("/devices/new", app.handleDeviceNew).Methods(http.MethodGet)
	router.HandleFunc("/devices/add", app.handleDeviceAdd).Methods(http.MethodPost)
	router.HandleFunc("/devices", app.handleDevices).Methods(http.MethodGet)
	return router
}
