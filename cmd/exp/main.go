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
	// if err := inventoryService.VendorRepo.New("fortigate"); err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println("successfully created a new vendor")
	// vendor, err := inventoryService.VendorRepo.GetByID(7)
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(vendor.Name)
	vendor, err := inventoryService.VendorRepo.GetByName("fortigate")
	if err != nil {
		log.Println(err)
	}
	fmt.Println(vendor.ID)
	vendor.Name = "updated fortigate"
	if err := inventoryService.VendorRepo.Update(vendor); err != nil {
		log.Println(err)
	}
	vendors, err := inventoryService.VendorRepo.GetAll()
	if err != nil {
		log.Println(err)
	}
	for _, v := range vendors {
		fmt.Println(v.ID, v.Name)
	}
	if err := inventoryService.VendorRepo.Delete(vendor.ID); err != nil {
		log.Println(err)
	}
	vendors, err = inventoryService.VendorRepo.GetAll()
	if err != nil {
		log.Println(err)
	}
	for _, v := range vendors {
		fmt.Println(v.ID, v.Name)
	}
}
