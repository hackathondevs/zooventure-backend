package repository

import (
	"context"

	"github.com/gofiber/fiber/v2"
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
	GetAllByUserID(ctx context.Context, userID int64) ([]model.Campaign, error)
	GetWithSubmission(ctx context.Context, id, userID int64) (model.Campaign, error)
	Create(ctx context.Context, campaign model.Campaign) error
	Update(ctx context.Context, campaign model.Campaign) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (model.Campaign, error)
	GetAll(ctx context.Context) ([]model.Campaign, error)
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

func (q *campaignQueryer) GetAllByUserID(ctx context.Context, userID int64) ([]model.Campaign, error) {
	var campaigns []model.Campaign
	if err := sqlx.SelectContext(ctx, q.ext, &campaigns, qGetAllCampaignByUserID, userID); err != nil {
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

func (q *campaignQueryer) Create(ctx context.Context, campaign model.Campaign) error {
	query, args, err := sqlx.Named(qCreateCampaign, fiber.Map{
		"Picture":     campaign.Picture,
		"Title":       campaign.Title,
		"Description": campaign.Description,
		"Reward":      campaign.Reward,
	})
	if err != nil {
		return err
	}
	_, err = q.ext.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return err
}

func (q *campaignQueryer) Update(ctx context.Context, campaign model.Campaign) error {
	query, args, err := sqlx.Named(qUpdateCampaign, fiber.Map{
		"ID":          campaign.ID,
		"Picture":     campaign.Picture,
		"Title":       campaign.Title,
		"Description": campaign.Description,
		"Reward":      campaign.Reward,
	})
	if err != nil {
		return err
	}
	_, err = q.ext.ExecContext(ctx, query, args...)
	return err
}

func (q *campaignQueryer) Delete(ctx context.Context, id int64) error {
	_, err := q.ext.ExecContext(ctx, qDeleteCampaign, id)
	if err != nil {
		return err
	}

	return nil
}

func (q *campaignQueryer) GetByID(ctx context.Context, id int64) (model.Campaign, error) {
	var campaign model.Campaign
	if err := sqlx.GetContext(ctx, q.ext, &campaign, qGetCampaignByID, id); err != nil {
		return model.Campaign{}, err
	}
	return campaign, nil
}

func (q *campaignQueryer) GetAll(ctx context.Context) ([]model.Campaign, error) {
	var campaigns []model.Campaign
	if err := sqlx.SelectContext(ctx, q.ext, &campaigns, qGetAllCampaign); err != nil {
		return nil, err
	}
	return campaigns, nil
}
