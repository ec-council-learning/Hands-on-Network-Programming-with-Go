package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/codered-by-ec-council/Hands-on-Network-Programming-with-Go/pkg/inventory"
	"github.com/codered-by-ec-council/Hands-on-Network-Programming-with-Go/pkg/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
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
	inventoryService := inventory.NewService(dbpool)
	testModel := models.Model{
		Name:   "test-model",
		Vendor: models.Vendor{ID: 1},
	}
	if err := inventoryService.ModelRepo.New(testModel); err != nil {
		log.Println(err)
	}
	fmt.Println("successfully created a new model")
	fmt.Println("getting model by ID")
	model, err := inventoryService.ModelRepo.GetByID(3)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(model)
	fmt.Println("getting model by Name")
	model, err = inventoryService.ModelRepo.GetByName("test-model")
	if err != nil {
		log.Println(err)
	}
	fmt.Println(model.ID)
	model.Name = "updated-test-model"
	fmt.Println("updating model name")
	if err := inventoryService.ModelRepo.Update(model); err != nil {
		log.Println(err)
	}
	fmt.Println("getting all models")
	models, err := inventoryService.ModelRepo.GetAll()
	if err != nil {
		log.Println(err)
	}
	for _, m := range models {
		fmt.Println(m.ID, m.Name, m.Vendor.ID, m.Vendor.Name)
	}
	fmt.Println("delting test model")
	if err := inventoryService.ModelRepo.Delete(model.ID); err != nil {
		log.Println(err)
	}
	fmt.Println("getting all models")
	models, err = inventoryService.ModelRepo.GetAll()
	if err != nil {
		log.Println(err)
	}
	for _, m := range models {
		fmt.Println(m.ID, m.Name, m.Vendor.ID, m.Vendor.Name)
	}
}
