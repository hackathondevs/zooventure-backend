package usecase

import (
	"context"
	"io"
	"path"
	"strings"

	"github.com/mirzahilmi/hackathon/internal/model"
	"github.com/mirzahilmi/hackathon/internal/pkg/gemini"
	"github.com/mirzahilmi/hackathon/internal/repository"
)

type AnimalUsecaseItf interface {
	PredictAnimal(ctx context.Context, raw *model.PredictAnimalRequest) (model.Animal, error)
}

type animalUsecase struct {
	userRepo      repository.UserRepositoryItf
	enclosureRepo repository.EnclosureRepositoryItf
	geminiModel   *gemini.GeminiAI
}

func NewAnimalUsecase(
	userRepo repository.UserRepositoryItf,
	enclosureRepo repository.EnclosureRepositoryItf,
) AnimalUsecaseItf {
	return &animalUsecase{userRepo, enclosureRepo, gemini.NewGeminiAI()}
}

func (u *animalUsecase) PredictAnimal(ctx context.Context, raw *model.PredictAnimalRequest) (model.Animal, error) {
	file, err := raw.Picture.Open()
	defer file.Close()
	if err != nil {
		return model.Animal{}, err
	}
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return model.Animal{}, err
	}
	prediction := u.geminiModel.PredictImageAnimal(ctx, fileBytes, strings.Replace(path.Ext(raw.Picture.Filename), ".", "", -1))
	if prediction.Name == "not animal" {
		return prediction, nil
	}
	enclosureRepo, err := u.enclosureRepo.NewClient(false, nil)
	if err != nil {
		return model.Animal{}, err
	}
	distance, err := enclosureRepo.FetchClosest(ctx, raw.Lat, raw.Long)
	if err != nil {
		return model.Animal{}, err
	}
	if distance < 16 {
		return prediction, nil
	}
	prediction.GotBonus = true
	userClient, err := u.userRepo.NewClient(false, nil)
	if err != nil {
		return model.Animal{}, err
	}
	if err := userClient.UpdateBalance(ctx, ctx.Value(ClientID).(int64), 100); err != nil {
		return model.Animal{}, err
	}
	return prediction, nil
}
