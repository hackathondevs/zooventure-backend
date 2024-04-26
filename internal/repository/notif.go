package repository

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type NotifRepositoryItf interface {
	NewClient(withTx bool, tx sqlx.ExtContext) (NotifQueryerItf, error)
}

type NotifRepository struct {
	db *sqlx.DB
}

func NewNotifRepository(db *sqlx.DB) NotifRepositoryItf {
	return &NotifRepository{db}
}

type NotifQueryerItf interface {
	txCompat
	Create(ctx context.Context, userID int64, txt string) error
	Fetch(ctx context.Context, id int64) ([]string, error)
}

type notifQueryer struct {
	ext sqlx.ExtContext
}

func (r *NotifRepository) NewClient(withTx bool, tx sqlx.ExtContext) (NotifQueryerItf, error) {
	if withTx {
		if tx != nil {
			return &notifQueryer{tx}, nil
		}
		tx, err := r.db.Beginx()
		if err != nil {
			return nil, err
		}
		return &notifQueryer{tx}, nil
	}
	return &notifQueryer{r.db}, nil
}

func (q *notifQueryer) Commit() error {
	tx, ok := q.ext.(*sqlx.Tx)
	if !ok {
		return ErrNotTransaction
	}
	return tx.Commit()
}

func (q *notifQueryer) Rollback() error {
	tx, ok := q.ext.(*sqlx.Tx)
	if !ok {
		return ErrNotTransaction
	}
	return tx.Rollback()
}

func (q *notifQueryer) Ext() sqlx.ExtContext {
	return q.ext
}

func (q *notifQueryer) Fetch(ctx context.Context, id int64) ([]string, error) {
	rows, err := q.ext.QueryxContext(ctx, qFetchNotifications, id)
	if err != nil {
		return nil, err
	}
	var notifs []string
	if err := rows.Scan(&notifs); err != nil {
		return nil, err
	}
	return notifs, err
}

func (q *notifQueryer) Create(ctx context.Context, userID int64, txt string) error {
	query, args, err := sqlx.Named(qCreateNotif, fiber.Map{"UserID": userID, "Text": txt})
	if err != nil {
		return err
	}
	if _, err := q.ext.ExecContext(ctx, query, args...); err != nil {
		return err
	}
	return nil
}
