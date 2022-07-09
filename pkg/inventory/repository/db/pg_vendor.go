package db

import (
	"context"
	"time"

	"github.com/codered-by-ec-council/Hands-on-Network-Programming-with-Go/pkg/models"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

const (
	newVendorSQL       = `INSERT INTO vendors (name) VALUES ($1)`
	getVendorByIDSQL   = `SELECT id, name FROM vendors WHERE id = $1`
	getVendorByNameSQL = `SELECT id, name FROM vendors WHERE name = $1`
	getVendorsSQL      = `SELECT id, name FROM vendors`
	updateVendorSQL    = `UPDATE vendors set name = $1 WHERE id = $2`
	deleteVendorSQL    = `DELETE FROM vendors WHERE id = $1`
)

type PGVendor struct {
	DBPool *pgxpool.Pool
}

func (pg *PGVendor) New(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := pg.DBPool.Exec(ctx, newVendorSQL, name)
	if err != nil {
		return errors.Wrap(err, "new vendor failed")
	}
	return nil
}

func (pg *PGVendor) GetByID(id int) (models.Vendor, error) {
	var vendor models.Vendor
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := pg.DBPool.QueryRow(ctx, getVendorByIDSQL, id).Scan(&vendor.ID, &vendor.Name); err != nil {
		return vendor, errors.Wrap(err, "GetByID QueryRow failed")
	}
	return vendor, nil
}

func (pg *PGVendor) GetByName(name string) (models.Vendor, error) {
	var vendor models.Vendor
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := pg.DBPool.QueryRow(ctx, getVendorByNameSQL, name).Scan(&vendor.ID, &vendor.Name); err != nil {
		return vendor, errors.Wrap(err, "GetByID QueryRow failed")
	}
	return vendor, nil
}

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

func (pg *PGVendor) Update(vendor models.Vendor) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := pg.DBPool.Exec(ctx, updateVendorSQL, vendor.Name, vendor.ID)
	if err != nil {
		return errors.Wrap(err, "Vendor update failed")
	}
	return nil
}

func (pg *PGVendor) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := pg.DBPool.Exec(ctx, deleteVendorSQL, id)
	if err != nil {
		return errors.Wrap(err, "Delete vendor failed")
	}
	return nil
}
