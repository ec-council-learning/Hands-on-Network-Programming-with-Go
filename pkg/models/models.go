package models

// A Vendor is a company that makes hardware devcies.
type Vendor struct {
	ID   int
	Name string
}

// A Model is a type of device from a vendor.
type Model struct {
	ID     int
	Name   string
	Vendor Vendor
}

// A Device is a specific instance of a vendor's model.
type Device struct {
	ID       int
	Hostname string
	IPv4     string
	Model    Model
}
