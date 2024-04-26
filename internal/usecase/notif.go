package usecase

import (
	"context"

	"github.com/mirzahilmi/hackathon/internal/repository"
	"github.com/sirupsen/logrus"
)

type NotifUsecaseItf interface {
	Create(ctx context.Context, txt string) error
	Fetch(ctx context.Context) ([]string, error)
}

type notifUsecase struct {
	notifRepo repository.NotifRepositoryItf
	log       *logrus.Logger
}

func NewNotifUsecase(
	notifRepo repository.NotifRepositoryItf,
	log *logrus.Logger,
) NotifUsecaseItf {
	return &notifUsecase{notifRepo, log}
}

func (u *notifUsecase) Fetch(ctx context.Context) ([]string, error) {
	client, err := u.notifRepo.NewClient(false, nil)
	if err != nil {
		return nil, err
	}
	notifs, err := client.Fetch(ctx, ctx.Value(ClientID).(int64))
	if err != nil {
		return nil, err
	}
	return notifs, nil
}

func (u *notifUsecase) Create(ctx context.Context, txt string) error {
	client, err := u.notifRepo.NewClient(false, nil)
	if err != nil {
		return err
	}
	if err := client.Create(ctx, ctx.Value(ClientID).(int64), txt); err != nil {
		return err
	}
	return nil
}
