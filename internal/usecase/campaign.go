package usecase

import (
	"context"

	"github.com/mirzahilmi/hackathon/internal/model"
	"github.com/mirzahilmi/hackathon/internal/repository"
	"github.com/sirupsen/logrus"
)

type CampaignUsecaseItf interface {
	FetchAll(ctx context.Context) ([]model.Campaign, error)
	GetWithSubmission(ctx context.Context, id int64) (model.Campaign, error)
}

type campaignUsecase struct {
	campaignRepo repository.CampaignRepositoryItf
	log          *logrus.Logger
}

func NewCampaignUsecase(
	campaignRepo repository.CampaignRepositoryItf,
	log *logrus.Logger,
) CampaignUsecaseItf {
	return &campaignUsecase{campaignRepo, log}
}

func (u *campaignUsecase) FetchAll(ctx context.Context) ([]model.Campaign, error) {
	client, err := u.campaignRepo.NewClient(false, nil)
	if err != nil {
		return nil, err
	}
	campaigns, err := client.GetAll(ctx, ctx.Value(ClientID).(int64))
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
