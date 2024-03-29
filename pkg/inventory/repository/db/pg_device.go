package db

import (
	"context"
	"time"

	"github.com/codered-by-ec-council/Hands-on-Network-Programming-with-Go/pkg/models"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

const newDeviceSQL = `-- name: newDevice: one
INSERT INTO devices (
	hostname,
	ipv4,
	model_id
) VALUES (
	$1,
	$2,
	$3
)`

const getDeviceByID = `-- name: getDeviceByID: one
SELECT
	devices.id,
	hostname,
	ipv4,
	vendors.name AS vendor,
	models.name AS model
FROM devices
JOIN models
	ON devices.model_id = models.id
JOIN vendors
	ON models.vendor_id = vendors.id
WHERE devices.id = $1`

const getDeviceByHostname = `-- name: getDeviceByHostname: one
SELECT
	devices.id,
	hostname,
	ipv4,
	vendors.name AS vendor,
	models.name AS model
FROM devices
JOIN models
	ON devices.model_id = models.id
JOIN vendors
	ON models.vendor_id = vendors.id
WHERE hostname = $1`

const getAllSQL = `-- name: getAll: many
SELECT
	devices.id,
	hostname,
	ipv4,
	vendors.name AS vendor,
	models.name AS model
FROM devices
JOIN models
	ON devices.model_id = models.id
JOIN vendors
	ON models.vendor_id = vendors.id`

const updateDeviceSQL = `-- name: updateDevice: one
UPDATE devices
SET
	hostname = $1,
	ipv4 = $2,
	model_id = $3
WHERE id = $4`

const deleteDeviceSQL = `-- name: deleteDevice: one
DELETE FROM devices
WHERE id = $1`

// PGDevice holds the DB connection pool.
type PGDevice struct {
	DBPool *pgxpool.Pool
}

// New adds a device to the table.
func (pg *PGDevice) New(device models.Device) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := pg.DBPool.Exec(ctx, newDeviceSQL, device.Hostname, device.IPv4, device.Model.ID)
	if err != nil {
		return errors.Wrap(err, "New device failed")
	}
	return nil
}

// GetByID returns a specific device matching the provided PK.
func (pg *PGDevice) GetByID(id int) (models.Device, error) {
	var device models.Device
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	var ip pgtype.Inet
	if err := pg.DBPool.QueryRow(ctx, getDeviceByID, id).Scan(
		&device.ID,
		&device.Hostname,
		&ip,
		&device.Model.Vendor.Name,
		&device.Model.Name,
	); err != nil {
		return device, errors.Wrap(err, "GetByID Scan failed")
	}
	device.IPv4 = ip.IPNet.IP.String()
	return device, nil
}

// GetByHostname returns a specific device matching the provided hostname.
func (pg *PGDevice) GetByHostname(name string) (models.Device, error) {
	var device models.Device
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	var ip pgtype.Inet
	if err := pg.DBPool.QueryRow(ctx, getDeviceByHostname, name).Scan(
		&device.ID,
		&device.Hostname,
		&ip,
		&device.Model.Vendor.Name,
		&device.Model.Name,
	); err != nil {
		return device, errors.Wrap(err, "GetByHostname Scan failed")
	}
	device.IPv4 = ip.IPNet.IP.String()
	return device, nil
}

// GetAll returns all devices in the table.
func (pg *PGDevice) GetAll() ([]models.Device, error) {
	var devices []models.Device
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	rows, err := pg.DBPool.Query(ctx, getAllSQL)
	if err != nil {
		return devices, errors.Wrap(err, "GetAll query failed")
	}
	defer rows.Close()
	for rows.Next() {
		var dev models.Device
		var ip pgtype.Inet
		if err := rows.Scan(
			&dev.ID,
			&dev.Hostname,
			&ip,
			&dev.Model.Vendor.Name,
			&dev.Model.Name,
		); err != nil {
			return devices, errors.Wrap(err, "GetByName Scan failed")
		}
		dev.IPv4 = ip.IPNet.IP.String()
		devices = append(devices, dev)
	}
	if rows.Err() != nil {
		return devices, errors.Wrap(err, "GetByName rows failed")
	}
	return devices, nil
}

// Update modifies an existing device.
func (pg *PGDevice) Update(device models.Device) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := pg.DBPool.Exec(ctx, updateDeviceSQL, device.Hostname, device.IPv4, device.Model.ID, device.ID)
	if err != nil {
		return errors.Wrap(err, "device update failed")
	}
	return nil
}

// Delete removes a device by it's PK.
func (pg *PGDevice) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := pg.DBPool.Exec(ctx, deleteDeviceSQL, id)
	if err != nil {
		return errors.Wrap(err, "device delete failed")
	}
	return nil
}
