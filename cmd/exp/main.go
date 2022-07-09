package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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
}
