package repository

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/mirzahilmi/hackathon/internal/model"
)

type ReportRepositoryItf interface {
	NewClient(withTx bool, tx sqlx.ExtContext) (ReportQueryerItf, error)
}

type ReportRepository struct {
	db *sqlx.DB
}

func NewReportRepository(db *sqlx.DB) ReportRepositoryItf {
	return &ReportRepository{db}
}

type ReportQueryerItf interface {
	txCompat
	CreateReport(ctx context.Context, report *model.ReportResource) error
	GetReports(ctx context.Context) ([]model.ReportResource, error)
	UpdateActionReport(ctx context.Context, report *model.ReportResource) error
	GetReportByID(ctx context.Context, id int64) (model.ReportResource, error)
}

type reportQueryer struct {
	ext sqlx.ExtContext
}

func (r *ReportRepository) NewClient(withTx bool, tx sqlx.ExtContext) (ReportQueryerItf, error) {
	if withTx {
		if tx != nil {
			return &reportQueryer{tx}, nil
		}
		tx, err := r.db.Beginx()
		if err != nil {
			return nil, err
		}
		return &reportQueryer{tx}, nil
	}
	return &reportQueryer{r.db}, nil
}

func (q *reportQueryer) Commit() error {
	tx, ok := q.ext.(*sqlx.Tx)
	if !ok {
		return ErrNotTransaction
	}
	return tx.Commit()
}

func (q *reportQueryer) Rollback() error {
	tx, ok := q.ext.(*sqlx.Tx)
	if !ok {
		return ErrNotTransaction
	}
	return tx.Rollback()
}

func (q *reportQueryer) Ext() sqlx.ExtContext {
	return q.ext
}

func (q *reportQueryer) CreateReport(ctx context.Context, report *model.ReportResource) error {
	query, args, err := sqlx.Named(qCreateReport, fiber.Map{
		"Picture":     report.Picture,
		"Description": report.Description,
		"Location":    report.Location,
	})
	log.Println(query)
	log.Println(args)
	if err != nil {
		return err
	}
	_, err = q.ext.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (q *reportQueryer) GetReports(ctx context.Context) ([]model.ReportResource, error) {
	var reports []model.ReportResource
	if err := sqlx.SelectContext(ctx, q.ext, &reports, qGetReports); err != nil {
		return nil, err
	}
	return reports, nil
}

func (q *reportQueryer) UpdateActionReport(ctx context.Context, report *model.ReportResource) error {
	query, args, err := sqlx.Named(qUpdateReport, fiber.Map{
		"Action": report.Action,
		"ID":     report.ID,
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

func (q *reportQueryer) GetReportByID(ctx context.Context, id int64) (model.ReportResource, error) {
	var report model.ReportResource
	if err := sqlx.GetContext(ctx, q.ext, &report, qGetReportByID, id); err != nil {
		return model.ReportResource{}, err
	}
	return report, nil
}
