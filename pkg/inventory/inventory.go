package inventory

import (
	"github.com/codered-by-ec-council/Hands-on-Network-Programming-with-Go/pkg/inventory/repository"
	"github.com/codered-by-ec-council/Hands-on-Network-Programming-with-Go/pkg/inventory/repository/db"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Service provides the fields to interact with the inventory
// persistance layer.
type Service struct {
	DBPool     *pgxpool.Pool
	VendorRepo repository.Vendor
	ModelRepo  repository.Model
	DeviceRepo repository.Device
}

// NewService spings up an inventory service.
func NewService(dbpool *pgxpool.Pool) *Service {
	return &Service{
		DBPool:     dbpool,
		VendorRepo: &db.PGVendor{DBPool: dbpool},
		ModelRepo:  &db.PGModel{DBPool: dbpool},
		DeviceRepo: &db.PGDevice{DBPool: dbpool},
	}
}
