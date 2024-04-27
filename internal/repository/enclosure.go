package repository

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type EnclosureRepositoryItf interface {
	NewClient(withTx bool, tx sqlx.ExtContext) (EnclosureQueryerItf, error)
}

type EnclosureRepository struct {
	db *sqlx.DB
}

func NewEnclosureRepository(db *sqlx.DB) EnclosureRepositoryItf {
	return &EnclosureRepository{db}
}

type EnclosureQueryerItf interface {
	txCompat
	FetchClosest(ctx context.Context, lat, long float64) (float64, error)
}

type enclosureQueryer struct {
	ext sqlx.ExtContext
}

func (r *EnclosureRepository) NewClient(withTx bool, tx sqlx.ExtContext) (EnclosureQueryerItf, error) {
	if withTx {
		if tx != nil {
			return &enclosureQueryer{tx}, nil
		}
		tx, err := r.db.Beginx()
		if err != nil {
			return nil, err
		}
		return &enclosureQueryer{tx}, nil
	}
	return &enclosureQueryer{r.db}, nil
}

func (q *enclosureQueryer) Commit() error {
	tx, ok := q.ext.(*sqlx.Tx)
	if !ok {
		return ErrNotTransaction
	}
	return tx.Commit()
}

func (q *enclosureQueryer) Rollback() error {
	tx, ok := q.ext.(*sqlx.Tx)
	if !ok {
		return ErrNotTransaction
	}
	return tx.Rollback()
}

func (q *enclosureQueryer) Ext() sqlx.ExtContext {
	return q.ext
}

func (q *enclosureQueryer) FetchClosest(ctx context.Context, lat float64, long float64) (float64, error) {
	query, args, err := sqlx.Named(qFetchDistanceToEnclosure, fiber.Map{"Latitude": lat, "Longitude": long})
	if err != nil {
		return 0, err
	}
	row := q.ext.QueryRowxContext(ctx, query, args...)
	var distance float64
	if err := row.Scan(&distance); err != nil {
		return 0, err
	}
	return distance, nil
}
