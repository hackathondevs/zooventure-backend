package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/mirzahilmi/hackathon/internal/model"
)

type CampaignRepositoryItf interface {
	NewClient(withTx bool, tx sqlx.ExtContext) (CampaignQueryerItf, error)
}

type CampaignRepository struct {
	db *sqlx.DB
}

func NewCampaignRepository(db *sqlx.DB) CampaignRepositoryItf {
	return &CampaignRepository{db}
}

type CampaignQueryerItf interface {
	txCompat
	GetAll(ctx context.Context, userID int64) ([]model.Campaign, error)
	GetWithSubmission(ctx context.Context, id, userID int64) (model.Campaign, error)
}

type campaignQueryer struct {
	ext sqlx.ExtContext
}

func (r *CampaignRepository) NewClient(withTx bool, tx sqlx.ExtContext) (CampaignQueryerItf, error) {
	if withTx {
		if tx != nil {
			return &campaignQueryer{tx}, nil
		}
		tx, err := r.db.Beginx()
		if err != nil {
			return nil, err
		}
		return &campaignQueryer{tx}, nil
	}
	return &campaignQueryer{r.db}, nil
}

func (q *campaignQueryer) Commit() error {
	tx, ok := q.ext.(*sqlx.Tx)
	if !ok {
		return ErrNotTransaction
	}
	return tx.Commit()
}

func (q *campaignQueryer) Rollback() error {
	tx, ok := q.ext.(*sqlx.Tx)
	if !ok {
		return ErrNotTransaction
	}
	return tx.Rollback()
}

func (q *campaignQueryer) Ext() sqlx.ExtContext {
	return q.ext
}

func (q *campaignQueryer) GetAll(ctx context.Context, userID int64) ([]model.Campaign, error) {
	var campaigns []model.Campaign
	if err := sqlx.SelectContext(ctx, q.ext, &campaigns, qGetAllCampaign, userID); err != nil {
		return nil, err
	}
	return campaigns, nil
}

func (q *campaignQueryer) GetWithSubmission(ctx context.Context, id, userID int64) (model.Campaign, error) {
	var campaign model.Campaign
	if err := sqlx.GetContext(ctx, q.ext, &campaign, qGetCampaignWithSubmission, id, userID); err != nil {
		return model.Campaign{}, err
	}
	return campaign, nil
}
