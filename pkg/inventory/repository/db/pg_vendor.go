package db

import (
	"context"
	"time"

	"github.com/codered-by-ec-council/Hands-on-Network-Programming-with-Go/pkg/models"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

const newVendorSQL = `-- name: newVendor: one
INSERT INTO vendors (
	name
) VALUES (
	$1
)`

const getVendorByIDSQL = `-- name: getVendorByID: one
SELECT
	id,
	name
FROM vendors WHERE id = $1`

const getVendorByNameSQL = `-- name: getVendorByName: one
SELECT
	id,
	name
FROM vendors
WHERE name = $1`

const getVendorsSQL = `-- name: getVendors: many
SELECT
	id,
	name
FROM vendors`

const updateVendorSQL = `-- name: updateVendor: one
UPDATE vendors
set
	name = $1
WHERE id = $2`

const deleteVendorSQL = `-- name: deleteVendor: one
DELETE FROM vendors
WHERE id = $1`

// PGVendor holds the DB connection pool.
type PGVendor struct {
	DBPool *pgxpool.Pool
}

// New adds a vendor to the table.
func (pg *PGVendor) New(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := pg.DBPool.Exec(ctx, newVendorSQL, name)
	if err != nil {
		return errors.Wrap(err, "new vendor failed")
	}
	return nil
}

// GetByID returns a specific vendor matching the provided PK.
func (pg *PGVendor) GetByID(id int) (models.Vendor, error) {
	var vendor models.Vendor
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := pg.DBPool.QueryRow(ctx, getVendorByIDSQL, id).Scan(&vendor.ID, &vendor.Name); err != nil {
		return vendor, errors.Wrap(err, "GetByID QueryRow failed")
	}
	return vendor, nil
}

// GetByName returns a specific vendor matching the provided name.
func (pg *PGVendor) GetByName(name string) (models.Vendor, error) {
	var vendor models.Vendor
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := pg.DBPool.QueryRow(ctx, getVendorByNameSQL, name).Scan(&vendor.ID, &vendor.Name); err != nil {
		return vendor, errors.Wrap(err, "GetByID QueryRow failed")
	}
	return vendor, nil
}

// GetAll returns all vendors in the table.
func (pg *PGVendor) GetAll() ([]models.Vendor, error) {
	var vendors []models.Vendor
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	rows, err := pg.DBPool.Query(ctx, getVendorsSQL)
	if err != nil {
		return vendors, errors.Wrap(err, "GetAll query failed")
	}
	defer rows.Close()
	for rows.Next() {
		var vendor models.Vendor
		if err := rows.Scan(&vendor.ID, &vendor.Name); err != nil {
			return vendors, errors.Wrap(err, "GetAll scan failed")
		}
		vendors = append(vendors, vendor)
	}
	if rows.Err() != nil {
		return vendors, errors.Wrap(err, "GetAll scan failed")
	}
	return vendors, nil
}

// Update modifies an existing vendor.
func (pg *PGVendor) Update(vendor models.Vendor) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := pg.DBPool.Exec(ctx, updateVendorSQL, vendor.Name, vendor.ID)
	if err != nil {
		return errors.Wrap(err, "Vendor update failed")
	}
	return nil
}

// Delete removes a vendor by it's PK.
func (pg *PGVendor) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := pg.DBPool.Exec(ctx, deleteVendorSQL, id)
	if err != nil {
		return errors.Wrap(err, "Delete vendor failed")
	}
	return nil
}
