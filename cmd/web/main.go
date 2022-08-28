package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/codered-by-ec-council/Hands-on-Network-Programming-with-Go/pkg/inventory"
	"github.com/jackc/pgx/v4/pgxpool"
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
	log.Println("starting web server on", *websocket)
	if err := http.ListenAndServe(*websocket, app.routes()); err != nil {
		log.Println(err)
	}
}
