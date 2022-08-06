package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	websocket := flag.String("websocket", "localhost:8080", "socket on which to listen for incoming connections")
	http.HandleFunc("/", handleHome)
	log.Println("starting web server on", *websocket)
	http.ListenAndServe(*websocket, nil)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Hello, go web app!</h1>"))
}
