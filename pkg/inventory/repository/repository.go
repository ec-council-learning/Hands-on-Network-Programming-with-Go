package repository

import "github.com/codered-by-ec-council/Hands-on-Network-Programming-with-Go/pkg/models"

type Vendor interface {
	New(name string) error
	GetByID(id int) (models.Vendor, error)
	GetByName(name string) (models.Vendor, error)
	GetAll() ([]models.Vendor, error)
	Update(vendor models.Vendor) error
	Delete(id int) error
}

type Model interface {
	New(model models.Model) error
	GetByID(id int) (models.Model, error)
	GetByName(name string) (models.Model, error)
	GetAll() ([]models.Model, error)
	Update(vendor models.Model) error
	Delete(id int) error
}

type Device interface {
	New(device models.Device) error
	GetByID(id int) (models.Device, error)
	GetByHostname(hostname string) (models.Device, error)
	GetAll() ([]models.Device, error)
	Update(device models.Device) error
	// Delete(id int) error
}
