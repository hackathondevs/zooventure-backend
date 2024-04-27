package usecase

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/mirzahilmi/hackathon/internal/model"
	"github.com/mirzahilmi/hackathon/internal/repository"
	"github.com/sirupsen/logrus"
	storage_go "github.com/supabase-community/storage-go"
)

type CampaignUsecaseItf interface {
	FetchAll(ctx context.Context) ([]model.Campaign, error)
	GetWithSubmission(ctx context.Context, id int64) (model.Campaign, error)
	Create(ctx context.Context, req model.CampaignRequest) error
	Update(ctx context.Context, campaign model.CampaignRequest, id int64) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (model.Campaign, error)
	GetAll(ctx context.Context) ([]model.Campaign, error)
	SubmitSubmission(ctx context.Context, id int64, req model.CampaignSubmissionRequest) error
}

type campaignUsecase struct {
	campaignRepo repository.CampaignRepositoryItf
	log          *logrus.Logger
	supabase     *storage_go.Client
}

func NewCampaignUsecase(
	campaignRepo repository.CampaignRepositoryItf,
	log *logrus.Logger,
	supabase *storage_go.Client,
) CampaignUsecaseItf {
	return &campaignUsecase{campaignRepo, log, supabase}
}

func (u *campaignUsecase) FetchAll(ctx context.Context) ([]model.Campaign, error) {
	client, err := u.campaignRepo.NewClient(false, nil)
	if err != nil {
		return nil, err
	}
	campaigns, err := client.GetAllByUserID(ctx, ctx.Value(ClientID).(int64))
	if err != nil {
		return nil, err
	}
	return campaigns, nil
}

func (u *campaignUsecase) GetWithSubmission(ctx context.Context, id int64) (model.Campaign, error) {
	client, err := u.campaignRepo.NewClient(false, nil)
	if err != nil {
		return model.Campaign{}, err
	}
	campaign, err := client.GetWithSubmission(ctx, id, ctx.Value(ClientID).(int64))
	if err != nil {
		return model.Campaign{}, err
	}
	return campaign, nil
}

func (u *campaignUsecase) Create(ctx context.Context, req model.CampaignRequest) error {
	client, err := u.campaignRepo.NewClient(true, nil)
	if err != nil {
		return err
	}
	defer client.Rollback()

	req.Picture.Filename = fmt.Sprintf("%d_%s", time.Now().UnixMilli(), req.Picture.Filename)
	pict, err := req.Picture.Open()
	if err != nil {
		return err
	}

	_, err = u.supabase.UploadFile(os.Getenv("SUPABASE_BUCKET_ID"), req.Picture.Filename, pict)
	if err != nil {
		return err
	}

	pictUrl := u.supabase.GetPublicUrl(os.Getenv("SUPABASE_BUCKET_ID"), req.Picture.Filename)

	campaign := model.Campaign{
		Picture:     pictUrl.SignedURL,
		Title:       req.Title,
		Description: req.Description,
		Reward:      req.Reward,
	}

	err = client.Create(ctx, campaign)
	if err != nil {
		return err
	}
	return client.Commit()
}

func (u *campaignUsecase) Update(ctx context.Context, req model.CampaignRequest, id int64) error {
	client, err := u.campaignRepo.NewClient(true, nil)
	if err != nil {
		return err
	}
	defer client.Rollback()

	campaign, err := client.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if req.Picture != nil {
		req.Picture.Filename = fmt.Sprintf("%d_%s", time.Now().UnixMilli(), req.Picture.Filename)
		pict, err := req.Picture.Open()
		if err != nil {
			return err
		}

		_, err = u.supabase.UploadFile(os.Getenv("SUPABASE_BUCKET_ID"), req.Picture.Filename, pict)
		if err != nil {
			return err
		}

		pictUrl := u.supabase.GetPublicUrl(os.Getenv("SUPABASE_BUCKET_ID"), req.Picture.Filename)
		campaign.Picture = pictUrl.SignedURL
	}

	if req.Title != "" && req.Title != campaign.Title {
		campaign.Title = req.Title
	}

	if req.Description != "" && req.Description != campaign.Description {
		campaign.Description = req.Description
	}

	if req.Reward != 0 && req.Reward != campaign.Reward {
		campaign.Reward = req.Reward
	}

	err = client.Update(ctx, campaign)
	if err != nil {
		return err
	}
	return client.Commit()
}

func (u *campaignUsecase) Delete(ctx context.Context, id int64) error {
	client, err := u.campaignRepo.NewClient(true, nil)
	if err != nil {
		return err
	}
	defer client.Rollback()
	err = client.Delete(ctx, id)
	if err != nil {
		return err
	}
	return client.Commit()
}

func (u *campaignUsecase) GetByID(ctx context.Context, id int64) (model.Campaign, error) {
	client, err := u.campaignRepo.NewClient(false, nil)
	if err != nil {
		return model.Campaign{}, err
	}
	campaign, err := client.GetByID(ctx, id)
	if err != nil {
		return model.Campaign{}, ErrCampaignNotExist
	}
	return campaign, nil
}

func (u *campaignUsecase) GetAll(ctx context.Context) ([]model.Campaign, error) {
	client, err := u.campaignRepo.NewClient(false, nil)
	if err != nil {
		return nil, err
	}
	campaigns, err := client.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return campaigns, nil
}

func (u *campaignUsecase) SubmitSubmission(ctx context.Context, id int64, req model.CampaignSubmissionRequest) error {
	client, err := u.campaignRepo.NewClient(true, nil)
	if err != nil {
		return err
	}
	defer client.Rollback()

	submission := model.CampaignSubmission{
		CampaignID: id,
		UserID:     ctx.Value(ClientID).(int64),
		Submission: req.Submission,
	}

	err = client.CreateSubmission(ctx, submission)
	if err != nil {
		return err
	}
	return client.Commit()
}
