package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/mirzahilmi/hackathon/internal/model"
)

type AnimalRepositoryItf interface {
	NewClient(withTx bool, tx sqlx.ExtContext) (AnimalQueryerItf, error)
}

type AnimalRepository struct {
	db *sqlx.DB
}

type AnimalQueryerItf interface {
	txCompat
	Save(ctx context.Context, animal *model.Animal) (int64, error)
}

type animalQueryer struct {
	ext sqlx.ExtContext
}

func (r *AnimalRepository) NewClient(withTx bool, tx sqlx.ExtContext) (AnimalQueryerItf, error) {
	if withTx {
		if tx != nil {
			return &animalQueryer{tx}, nil
		}
		tx, err := r.db.Beginx()
		if err != nil {
			return nil, err
		}
		return &animalQueryer{tx}, nil
	}
	return &animalQueryer{r.db}, nil
}

func (q *animalQueryer) Commit() error {
	tx, ok := q.ext.(*sqlx.Tx)
	if !ok {
		return ErrNotTransaction
	}
	return tx.Commit()
}

func (q *animalQueryer) Rollback() error {
	tx, ok := q.ext.(*sqlx.Tx)
	if !ok {
		return ErrNotTransaction
	}
	return tx.Rollback()
}

func (q *animalQueryer) Ext() sqlx.ExtContext {
	return q.ext
}
func (q *animalQueryer) Save(ctx context.Context, animal *model.Animal) (int64, error) {
	return 0, nil
}
