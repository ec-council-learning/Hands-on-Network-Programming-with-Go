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
