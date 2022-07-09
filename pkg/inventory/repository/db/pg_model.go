package db

import (
	"context"
	"time"

	"github.com/codered-by-ec-council/Hands-on-Network-Programming-with-Go/pkg/models"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

const (
	newModelSQL     = `INSERT INTO models (name, vendor_id) VALUES ($1, $2)`
	getModelByIDSQL = `
		SELECT models.id, models.name, vendors.id, vendors.name
		FROM models
		JOIN vendors ON vendors.id = models.vendor_id
		WHERE models.id = $1`
	getModelByNameSQL = `
		SELECT models.id, models.name, vendors.id, vendors.name
		FROM models
		JOIN vendors ON vendors.id = models.vendor_id
		WHERE models.name = $1`
	getModelsSQL = `
		SELECT models.id, models.name, vendors.id, vendors.name
		FROM models
		JOIN vendors ON vendors.id = models.vendor_id`
	updateModelSQL = `UPDATE models set name = $1, vendor_id = $2 WHERE id = $3`
	deleteModelSQL = `DELETE FROM models WHERE id = $1`
)

type PGModel struct {
	DBPool *pgxpool.Pool
}

func (pg *PGModel) New(model models.Model) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := pg.DBPool.Exec(ctx, newModelSQL, model.Name, model.Vendor.ID)
	if err != nil {
		return errors.Wrap(err, "new model failed")
	}
	return nil
}

func (pg *PGModel) GetByID(id int) (models.Model, error) {
	var model models.Model
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := pg.DBPool.QueryRow(ctx, getModelByIDSQL, id).Scan(
		&model.ID,
		&model.Name,
		&model.Vendor.ID,
		&model.Vendor.Name,
	); err != nil {
		return model, errors.Wrap(err, "GetByID QueryRow failed")
	}
	return model, nil
}

func (pg *PGModel) GetByName(name string) (models.Model, error) {
	var model models.Model
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := pg.DBPool.QueryRow(ctx, getModelByNameSQL, name).Scan(
		&model.ID,
		&model.Name,
		&model.Vendor.ID,
		&model.Vendor.Name,
	); err != nil {
		return model, errors.Wrap(err, "GetByID QueryRow failed")
	}
	return model, nil
}

func (pg *PGModel) GetAll() ([]models.Model, error) {
	var mods []models.Model
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	rows, err := pg.DBPool.Query(ctx, getModelsSQL)
	if err != nil {
		return mods, errors.Wrap(err, "GetAll query failed")
	}
	defer rows.Close()
	for rows.Next() {
		var m models.Model
		if err := rows.Scan(
			&m.ID,
			&m.Name,
			&m.Vendor.ID,
			&m.Vendor.Name,
		); err != nil {
			return mods, errors.Wrap(err, "GetAll scan failed")
		}
		mods = append(mods, m)
	}
	if rows.Err() != nil {
		return mods, errors.Wrap(err, "GetAll scan failed")
	}
	return mods, nil
}

func (pg *PGModel) Update(model models.Model) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := pg.DBPool.Exec(ctx, updateModelSQL, model.Name, model.Vendor.ID, model.ID)
	if err != nil {
		return errors.Wrap(err, "Model update failed")
	}
	return nil
}

func (pg *PGModel) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := pg.DBPool.Exec(ctx, deleteModelSQL, id)
	if err != nil {
		return errors.Wrap(err, "Delete model failed")
	}
	return nil
}
