package models

type Vendor struct {
	ID   int
	Name string
}

type Model struct {
	ID     int
	Name   string
	Vendor Vendor
}

type Device struct {
	ID       int
	Hostname string
	IPv4     string
	Model    Model
}
