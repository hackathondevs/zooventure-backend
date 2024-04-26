package gemini

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
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

type AnimalPredictResp struct {
	Name            string   `json:"name"`
	LatinName       string   `json:"latin_name"`
	CountryOfOrigin string   `json:"country_of_origin"`
	Characteristics []string `json:"characteristics"`
	Category        string   `json:"category"`
	Lifespan        string   `json:"lifespan"`
	Funfact         string   `json:"funfact"`
}

func (g *GeminiAI) PredictImageAnimal(ctx context.Context, data []byte, typeFile string) any {
	resp, err := g.model.GenerateContent(ctx,
		genai.Text("can you describe this image based and show with format json that i perform"),
		genai.Text("the json using this format:"),
		genai.Text("{"),
		genai.Text("name : string,"),
		genai.Text("latin_name : string,"),
		genai.Text("country_of_origin : string,"),
		genai.Text("characteristics : array of string,"),
		genai.Text("category : string (Karnivora or Omnivora or Herbivora),"),
		genai.Text("lifespan : string,"),
		genai.Text("funfact : string (fun fact about the animal),"),
		genai.Text("}"),
		genai.Text("if the image not animal and not have any animal in the image please just return not animal string"),
		genai.Text("Please provide the response as json format (inside json)"),
		genai.ImageData("jpeg", data),
	)
	if err != nil {
		log.Fatal(err)
	}

	part := resp.Candidates[0].Content.Parts[0]

	jsonPart, err := json.Marshal(part)
	if err != nil {
		log.Fatal(err)
	}

	input := string(jsonPart)
	cleanedInput := strings.ReplaceAll(input, "\\n", "")
	cleanedInput = strings.ReplaceAll(cleanedInput, "\\\"", "\"")
	cleanedInput = cleanedInput[2 : len(cleanedInput)-1]

	var AnimalPredictResp AnimalPredictResp
	if err := json.Unmarshal([]byte(cleanedInput), &AnimalPredictResp); err != nil {
		log.Fatal(err)
	}

	return AnimalPredictResp
}
