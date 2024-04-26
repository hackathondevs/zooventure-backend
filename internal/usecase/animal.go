package usecase

import (
	"context"
	"io"
	"path"
	"strings"

	"github.com/mirzahilmi/hackathon/internal/model"
	"github.com/mirzahilmi/hackathon/internal/pkg/gemini"
)

type AnimalUsecaseItf interface {
	PredictAnimal(ctx context.Context, raw *model.PredictAnimalRequest) (model.Animal, error)
}

type animalUsecase struct {
	geminiModel *gemini.GeminiAI
}

func NewAnimalUsecase() AnimalUsecaseItf {
	geminiModel := gemini.NewGeminiAI()
	return &animalUsecase{geminiModel}
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
	predict := u.geminiModel.PredictImageAnimal(ctx, fileBytes, strings.Replace(path.Ext(raw.Picture.Filename), ".", "", -1))
	if predict.Name == "not animal" {
		return predict, nil
	}

	return predict, nil
}
