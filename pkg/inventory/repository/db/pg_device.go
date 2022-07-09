package db

import (
	"context"
	"time"

	"github.com/codered-by-ec-council/Hands-on-Network-Programming-with-Go/pkg/models"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

const (
	newDeviceSQL = `INSERT INTO devices (hostname, ipv4, model_id) VALUES ($1, $2, $3)`
)

type PGDevice struct {
	DBPool *pgxpool.Pool
}

func (pg *PGDevice) New(device models.Device) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := pg.DBPool.Exec(ctx, newDeviceSQL, device.Hostname, device.IPv4, device.Model.ID)
	if err != nil {
		return errors.Wrap(err, "New device failed")
	}
	return nil
}
