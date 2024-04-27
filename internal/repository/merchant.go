package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/mirzahilmi/hackathon/internal/model"
)

type MerchantRepositoryItf interface {
	NewClient(withTx bool, tx sqlx.ExtContext) (MerchantQueryerItf, error)
}

type MerchantRepository struct {
	db *sqlx.DB
}

func NewMerchantRepository(db *sqlx.DB) MerchantRepositoryItf {
	return &MerchantRepository{db}
}

type MerchantQueryerItf interface {
	txCompat
	GetByParam(ctx context.Context, param, value interface{}) (model.MerchantResource, error)
}

type merchantQueryer struct {
	ext sqlx.ExtContext
}

func (r *MerchantRepository) NewClient(withTx bool, tx sqlx.ExtContext) (MerchantQueryerItf, error) {
	if withTx {
		if tx != nil {
			return &merchantQueryer{tx}, nil
		}
		tx, err := r.db.Beginx()
		if err != nil {
			return nil, err
		}
		return &merchantQueryer{tx}, nil
	}
	return &merchantQueryer{r.db}, nil
}

func (q *merchantQueryer) Commit() error {
	tx, ok := q.ext.(*sqlx.Tx)
	if !ok {
		return ErrNotTransaction
	}
	return tx.Commit()
}

func (q *merchantQueryer) Rollback() error {
	tx, ok := q.ext.(*sqlx.Tx)
	if !ok {
		return ErrNotTransaction
	}
	return tx.Rollback()
}

func (q *merchantQueryer) Ext() sqlx.ExtContext {
	return q.ext
}

func (q *merchantQueryer) GetByParam(ctx context.Context, param, value interface{}) (model.MerchantResource, error) {
	row := q.ext.QueryRowxContext(ctx, fmt.Sprintf(qGetMerchantByParam, param), value)
	if err := row.Err(); err != nil {
		return model.MerchantResource{}, err
	}
	var merchant model.MerchantResource
	if err := row.StructScan(&merchant); err != nil {
		return model.MerchantResource{}, err
	}
	return merchant, nil
}
