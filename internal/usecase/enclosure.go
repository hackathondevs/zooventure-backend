package usecase

import (
	"context"

	"github.com/mirzahilmi/hackathon/internal/model"
	"github.com/mirzahilmi/hackathon/internal/repository"
	"github.com/sirupsen/logrus"
)

type EnclosureUsecaseItf interface {
	GetAll(ctx context.Context) ([]model.EnclosureResource, error)
}

type enclosureUsecase struct {
	enclosureRepo repository.EnclosureRepositoryItf
	log           *logrus.Logger
}

func NewEnclosureUsecase(
	enclosureRepo repository.EnclosureRepositoryItf,
	log *logrus.Logger,
) EnclosureUsecaseItf {
	return &enclosureUsecase{enclosureRepo, log}
}

func (u *enclosureUsecase) GetAll(ctx context.Context) ([]model.EnclosureResource, error) {
	client, err := u.enclosureRepo.NewClient(false, nil)
	if err != nil {
		return nil, err
	}
	enclosuresData, err := client.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	var enclosures []model.EnclosureResource
	for _, enclosure := range enclosuresData {
		enclosures = append(enclosures, enclosure.Clean())
	}
	return enclosures, nil
}
