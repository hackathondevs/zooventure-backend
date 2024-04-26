package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/mirzahilmi/hackathon/internal/model"
)

type VerifRepositoryItf interface {
	NewClient(withTx bool, tx sqlx.ExtContext) (VerifQueryerItf, error)
}

type VerifRepository struct {
	db *sqlx.DB
}

func NewVerificationRepository(db *sqlx.DB) VerifRepositoryItf {
	return &VerifRepository{db}
}

type VerifQueryerItf interface {
	txCompat
	Create(ctx context.Context, verif *model.UserVerificationResource) (int64, error)
	GetByIDAndToken(ctx context.Context, id int64, token string) (model.UserVerificationResource, error)
	UpdateSucceedStatus(ctx context.Context, id int64) error
}

type verifQueryer struct {
	ext sqlx.ExtContext
}

func (r *VerifRepository) NewClient(withTx bool, tx sqlx.ExtContext) (VerifQueryerItf, error) {
	if withTx {
		if tx != nil {
			return &verifQueryer{tx}, nil
		}
		tx, err := r.db.Beginx()
		if err != nil {
			return nil, err
		}
		return &verifQueryer{tx}, nil
	}
	return &verifQueryer{r.db}, nil
}

func (q *verifQueryer) Commit() error {
	tx, ok := q.ext.(*sqlx.Tx)
	if !ok {
		return ErrNotTransaction
	}
	return tx.Commit()
}

func (q *verifQueryer) Rollback() error {
	tx, ok := q.ext.(*sqlx.Tx)
	if !ok {
		return ErrNotTransaction
	}
	return tx.Rollback()
}

func (q *verifQueryer) Ext() sqlx.ExtContext {
	return q.ext
}

func (q *verifQueryer) Create(ctx context.Context, verif *model.UserVerificationResource) (int64, error) {
	query, args, err := sqlx.Named(qCreateUserVerification, verif)
	if err != nil {
		return 0, err
	}
	res, err := q.ext.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (q *verifQueryer) GetByIDAndToken(ctx context.Context, id int64, token string) (model.UserVerificationResource, error) {
	row := q.ext.QueryRowxContext(ctx, qGetUserVerificationByIDAndToken, id, token)
	if err := row.Err(); err != nil {
		return model.UserVerificationResource{}, err
	}
	var verif model.UserVerificationResource
	if err := row.StructScan(&verif); err != nil {
		return model.UserVerificationResource{}, err
	}
	return verif, nil
}

func (q *verifQueryer) UpdateSucceedStatus(ctx context.Context, id int64) error {
	if _, err := q.ext.ExecContext(ctx, qUpdateUserVerificationStatus, id); err != nil {
		return err
	}
	return nil
}
