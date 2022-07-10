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
	updatedDevice := models.Device{
		ID:       2,
		Hostname: "updated-testhostname",
		IPv4:     "42.42.42.42",
		Model:    models.Model{ID: 1},
	}
	if err := inventoryService.DeviceRepo.Update(updatedDevice); err != nil {
		log.Println(err)
	}
}
