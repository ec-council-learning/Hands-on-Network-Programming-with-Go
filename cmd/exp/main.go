package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/codered-by-ec-council/Hands-on-Network-Programming-with-Go/pkg/inventory"
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
	dev, err := inventoryService.DeviceRepo.GetByHostname("labsrx")
	if err != nil {
		log.Println(err)
	}
	fmt.Println(dev.ID, dev.Hostname, dev.Model.Vendor.Name, dev.Model.Name)
	dev, err = inventoryService.DeviceRepo.GetByID(dev.ID)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(dev.Hostname)
	devices, err := inventoryService.DeviceRepo.GetAll()
	if err != nil {
		log.Println(err)
	}
	for _, d := range devices {
		fmt.Println(d.Hostname)
	}
}
