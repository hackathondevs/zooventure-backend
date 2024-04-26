package usecase

import (
	"context"
	"io/ioutil"
	"mime/multipart"
	"path"
	"strings"

	"github.com/mirzahilmi/hackathon/internal/pkg/gemini"
)

type AnimalUsecaseItf interface {
	PredictAnimal(ctx context.Context, pict *multipart.FileHeader) any
}

type animalUsecase struct {
	geminiModel *gemini.GeminiAI
}

func NewAnimalUsecase() AnimalUsecaseItf {
	geminiModel := gemini.NewGeminiAI()

	return &animalUsecase{
		geminiModel: geminiModel,
	}
}

func (u *animalUsecase) PredictAnimal(ctx context.Context, pict *multipart.FileHeader) any {
	file, _ := pict.Open()
	fileBytes, _ := ioutil.ReadAll(file)
	defer file.Close()

	respChan := make(chan any)
	defer close(respChan)

	go func() {
		respChan <- u.geminiModel.PredictImageAnimal(ctx, fileBytes, strings.Replace(path.Ext(pict.Filename), ".", "", -1))
	}()

	resp := <-respChan

	return &resp
}
