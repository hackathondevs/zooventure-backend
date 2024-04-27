package gemini

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"

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
		genai.Text("Forget ANY CONTEXT before this chat, and listen carefully what i want you to do"),
		genai.Text("can you describe this animal picture, if the picture does not seem to be an animal picture, just fill out all the string field with \"not animal\""),
		genai.Text("use this JSON format that i provide below for you to search its information"),
		genai.Text(`{"name":"The animal name","latin":"The animal latin name","countryOfOrigin":"Countries, where, animal, origins, separated in, commas","characteristics": ["Animal","Unique","Characteristic"],"category": "Carnivore or Herbivore or Omnivore","lifespan":"animal lifespan","funfact": "animal's fun fact that not so many people know"}`),
		genai.Text("Please only provide me with the data in json format in one-line, DONT output anything else"),
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
	raw = strings.ReplaceAll(raw, "```json", "")
	raw = strings.ReplaceAll(raw, "```", "")
	var animal model.Animal
	if err := jsoniter.Unmarshal([]byte(raw), &animal); err != nil {
		return model.Animal{}, err
	}
	return animal, nil
}
