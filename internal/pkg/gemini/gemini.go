package gemini

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/google/generative-ai-go/genai"
	jsoniter "github.com/json-iterator/go"
	"github.com/mirzahilmi/hackathon/internal/model"
	"google.golang.org/api/option"
)

type GeminiAI struct {
	model *genai.GenerativeModel
}

func NewGeminiAI() *GeminiAI {
	geminiApiKey := os.Getenv("GEMINI_API_KEY")
	geminiTypeModel := os.Getenv("GEMINI_TYPE_MODEL")

	ctx := context.Background()

	client, err := genai.NewClient(ctx, option.WithAPIKey(geminiApiKey))
	if err != nil {
		log.Fatal(err)
	}

	model := client.GenerativeModel(geminiTypeModel)
	geminiAi := &GeminiAI{
		model: model,
	}

	return geminiAi
}

func (g *GeminiAI) PredictImageAnimal(ctx context.Context, data []byte, typeFile string) (model.Animal, error) {
	prompt, err := g.model.GenerateContent(ctx,
		genai.Text("can you describe this image based and show with format json that i perform and if the image is not animal please fill the json with not animal"),
		genai.Text("the json using this format:"),
		genai.Text("{"),
		genai.Text("name : string,"),
		genai.Text("latinName : string,"),
		genai.Text("countryOfOrigin : string,"),
		genai.Text("characteristics : array of string (max is 5 charater),"),
		genai.Text("category : string (Karnivora or Omnivora or Herbivora),"),
		genai.Text("lifespan : string,"),
		genai.Text("funfact : string (fun fact about the animal),"),
		genai.Text("}"),
		genai.Text("if the image is not animal and dont have any animal in the image please just fill the json stringfields with not animal"),
		genai.Text("Please only provide me with the data in json format one-line, DONT output anything else"),
		genai.ImageData("jpeg", data),
	)
	if err != nil {
		return model.Animal{}, err
	}
	resp := prompt.Candidates[0].Content.Parts[0]
	json, err := json.Marshal(resp)
	if err != nil {
		return model.Animal{}, err
	}
	raw, err := strconv.Unquote(string(json))
	if err != nil {
		return model.Animal{}, err
	}
	raw = raw[8 : len(raw)-3]
	var animal model.Animal
	if err := jsoniter.Unmarshal([]byte(raw), &animal); err != nil {
		return model.Animal{}, err
	}
	return animal, nil
}
