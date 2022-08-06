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
	w.Write([]byte(`
	<!doctype html>
	<html lang="en">
	  <head>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<title>Bootstrap demo</title>
		<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.2.0/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-gH2yIJqKdNHPEq0n4Mqa/HGKIhSkIHeL5AyhkYV8i59U5AR6csBvApHHNl/vI1Bx" crossorigin="anonymous">
	  </head>
	  <body>
	  <nav class="navbar navbar-expand-lg navbar-dark bg-dark">
	  <div class="container-fluid">
		<a class="navbar-brand" href="/">Gopher Engineering</a>
		<button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarSupportedContent" aria-controls="navbarSupportedContent" aria-expanded="false" aria-label="Toggle navigation">
		  <span class="navbar-toggler-icon"></span>
		</button>
		<div class="collapse navbar-collapse" id="navbarSupportedContent">
		  <ul class="navbar-nav me-auto mb-2 mb-lg-0">
			<li class="nav-item">
			  <a class="nav-link active" aria-current="page" href="/">Home</a>
			</li>
			<li class="nav-item dropdown">
			  <a class="nav-link dropdown-toggle" href="#" role="button" data-bs-toggle="dropdown" aria-expanded="false">
				Vendors
			  </a>
			  <ul class="dropdown-menu">
				<li><a class="dropdown-item" href="/vendors/new">New</a></li>
				<li><a class="dropdown-item" href="/vendors">List</a></li>
			  </ul>
			</li>
			<li class="nav-item dropdown">
			  <a class="nav-link dropdown-toggle" href="#" role="button" data-bs-toggle="dropdown" aria-expanded="false">
				Models
			  </a>
			  <ul class="dropdown-menu">
				<li><a class="dropdown-item" href="/models/new">New</a></li>
				<li><a class="dropdown-item" href="/models">List</a></li>
			  </ul>
			</li>
			<li class="nav-item dropdown">
			  <a class="nav-link dropdown-toggle" href="#" role="button" data-bs-toggle="dropdown" aria-expanded="false">
				Devices
			  </a>
			  <ul class="dropdown-menu">
				<li><a class="dropdown-item" href="/devices/new">New</a></li>
				<li><a class="dropdown-item" href="/devices">List</a></li>
			  </ul>
			</li>
		  </ul>
		  <form class="d-flex" role="search">
			<input class="form-control me-2" type="search" placeholder="Search" aria-label="Search">
			<button class="btn btn-outline-success" type="submit">Search</button>
		  </form>
		</div>
	  </div>
	</nav>
	<div class="container py-4">
		<div class="p-5 mb-4 bg-light rounded-3">
			<div class="container-fluid py-5">
				<h1 class="display-5 fw-bold">Gopher Engineering</h1>
				<p class="col-md-8 fs-4">Using a series of utilities, you can create this jumbotron, just like the one in previous versions of Bootstrap. Check out the examples below for how you can remix and restyle it to your liking.</p>
				<button class="btn btn-primary btn-lg" type="button">Docs</button>
			</div>
		</div>
	</div>
	<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.2.0/dist/js/bootstrap.bundle.min.js" integrity="sha384-A3rJD856KowSb7dwlZdYEkO39Gagi7vIsF0jrRAoQmDKKtQBHUuLZ9AsSv4jD4Xa" crossorigin="anonymous"></script>
	</body>
	</html>
	`))
}
