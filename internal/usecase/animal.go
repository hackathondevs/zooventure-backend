package usecase

import (
	"context"
	"io"

	"github.com/mirzahilmi/hackathon/internal/model"
	"github.com/mirzahilmi/hackathon/internal/pkg/gemini"
	"github.com/mirzahilmi/hackathon/internal/repository"
)

type AnimalUsecaseItf interface {
	PredictAnimal(ctx context.Context, data *model.PredictAnimalReq) (model.AnimalResource, error)
	FetchAll(ctx context.Context) ([]model.AnimalResource, error)
}

type animalUsecase struct {
	userRepo    repository.UserRepositoryItf
	animalRepo  repository.AnimalRepositoryItf
	geminiModel *gemini.GeminiAI
}

func NewAnimalUsecase(
	userRepo repository.UserRepositoryItf,
	animalRepo repository.AnimalRepositoryItf,
) AnimalUsecaseItf {
	return &animalUsecase{userRepo, animalRepo, gemini.NewGeminiAI()}
}

func (u *animalUsecase) PredictAnimal(ctx context.Context, data *model.PredictAnimalReq) (model.AnimalResource, error) {
	file, err := data.Picture.Open()
	if err != nil {
		return model.AnimalResource{}, err
	}
	defer file.Close()
	if err != nil {
		return model.AnimalResource{}, err
	}
	bytes, err := io.ReadAll(file)
	if err != nil {
		return model.AnimalResource{}, err
	}
	prediction, err := u.geminiModel.PredictImageAnimal(ctx, bytes)
	if err != nil {
		return model.AnimalResource{}, err
	}
	prediction = prediction[1:]
	if prediction == "notanimal" {
		return model.AnimalResource{}, err
	}
	animalClient, err := u.animalRepo.NewClient(false, nil)
	if err != nil {
		return model.AnimalResource{}, err
	}
	animal, err := animalClient.FetchTopRelated(ctx, "%"+prediction+"%", data.Lat, data.Long)
	if err != nil {
		return model.AnimalResource{}, err
	}

	return animal.Resource(), nil
}

// FetchAll implements AnimalUsecaseItf.
func (u *animalUsecase) FetchAll(ctx context.Context) ([]model.AnimalResource, error) {
	animalClient, err := u.animalRepo.NewClient(false, nil)
	if err != nil {
		return nil, err
	}
	animalsDat, err := animalClient.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	var animals []model.AnimalResource
	for _, animal := range animalsDat {
		animals = append(animals, animal.Resource())
	}
	return animals, nil
}
