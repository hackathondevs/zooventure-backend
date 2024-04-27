package repository

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/mirzahilmi/hackathon/internal/model"
)

type UserRepositoryItf interface {
	NewClient(withTx bool, tx sqlx.ExtContext) (UserQueryerItf, error)
}

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepositoryItf {
	return &UserRepository{db}
}

type UserQueryerItf interface {
	txCompat
	GetByParam(ctx context.Context, param, value interface{}) (model.UserResource, error)
	Create(ctx context.Context, user *model.UserResource) (int64, error)
	Update(ctx context.Context, user *model.UserResource) error
	UpdatePicture(ctx context.Context, id int64, url string) error
	UpdatePassword(ctx context.Context, id int64, passwd string) error
	UpdateBalance(ctx context.Context, id int64, value int) error
	UpdateActiveStatus(ctx context.Context, id int64) error
	CreateExchange(ctx context.Context, exchange *model.ExchangeResource) error
	GetExchanges(ctx context.Context, param, value interface{}) ([]model.ExchangeCleanResource, error)
}

type userQueryer struct {
	ext sqlx.ExtContext
}

func (r *UserRepository) NewClient(withTx bool, tx sqlx.ExtContext) (UserQueryerItf, error) {
	if withTx {
		if tx != nil {
			return &userQueryer{tx}, nil
		}
		tx, err := r.db.Beginx()
		if err != nil {
			return nil, err
		}
		return &userQueryer{tx}, nil
	}
	return &userQueryer{r.db}, nil
}

func (q *userQueryer) Commit() error {
	tx, ok := q.ext.(*sqlx.Tx)
	if !ok {
		return ErrNotTransaction
	}
	return tx.Commit()
}

func (q *userQueryer) Rollback() error {
	tx, ok := q.ext.(*sqlx.Tx)
	if !ok {
		return ErrNotTransaction
	}
	return tx.Rollback()
}

func (q *userQueryer) Ext() sqlx.ExtContext {
	return q.ext
}

func (q *userQueryer) GetByParam(ctx context.Context, param, value interface{}) (model.UserResource, error) {
	row := q.ext.QueryRowxContext(ctx, fmt.Sprintf(queryGetUserByParam, param), value)
	if err := row.Err(); err != nil {
		return model.UserResource{}, err
	}
	var user model.UserResource
	if err := row.StructScan(&user); err != nil {
		return model.UserResource{}, err
	}
	return user, nil
}

func (q *userQueryer) Create(ctx context.Context, user *model.UserResource) (int64, error) {
	query, args, err := sqlx.Named(qCreateUser, user)
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

func (q *userQueryer) Update(ctx context.Context, user *model.UserResource) error {
	query, args, err := sqlx.Named(qUpdateUserProfile, user)
	if err != nil {
		return err
	}
	if _, err := q.ext.ExecContext(ctx, query, args...); err != nil {
		return err
	}
	return nil
}

func (q *userQueryer) UpdatePicture(ctx context.Context, id int64, url string) error {
	query, args, err := sqlx.Named(fmt.Sprintf(qUpdateUserField, "ProfilePicture"), fiber.Map{"ID": id, "Value": url})
	if err != nil {
		return err
	}
	if _, err := q.ext.ExecContext(ctx, query, args...); err != nil {
		return err
	}
	return nil
}

func (q *userQueryer) UpdatePassword(ctx context.Context, id int64, passwd string) error {
	query, args, err := sqlx.Named(fmt.Sprintf(qUpdateUserField, "Password"), fiber.Map{"ID": id, "Value": passwd})
	if err != nil {
		return err
	}
	res, err := q.ext.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return ErrNoRowsAffected
	}
	return nil
}

func (q *userQueryer) UpdateActiveStatus(ctx context.Context, id int64) error {
	res, err := q.ext.ExecContext(ctx, qUpdateUserStatus, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return ErrNoRowsAffected
	}
	return nil
}

func (q *userQueryer) UpdateBalance(ctx context.Context, id int64, value int) error {
	query, args, err := sqlx.Named(qUpdateUserBalance, fiber.Map{"ID": id, "Value": value})
	if err != nil {
		return err
	}
	if _, err := q.ext.ExecContext(ctx, query, args...); err != nil {
		return err
	}
	return nil
}

func (q *userQueryer) CreateExchange(ctx context.Context, exchange *model.ExchangeResource) error {
	query, args, err := sqlx.Named(qCreateExchange,
		fiber.Map{
			"MerchantID": exchange.MerchantID,
			"UserID":     exchange.UserID,
			"Amount":     exchange.Amount,
			"Date":       exchange.Date,
			"Status":     exchange.Status,
		})
	if err != nil {
		return err
	}
	_, err = q.ext.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (q *userQueryer) GetExchanges(ctx context.Context, param, value interface{}) ([]model.ExchangeCleanResource, error) {
	rows, err := q.ext.QueryxContext(ctx, fmt.Sprintf(qGetExchanges, param), value)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var exchanges []model.ExchangeCleanResource
	for rows.Next() {
		var exchange model.ExchangeResource
		if err := rows.StructScan(&exchange); err != nil {
			return nil, err
		}
		exchanges = append(exchanges, exchange.Clean())
	}
	return exchanges, nil
}
