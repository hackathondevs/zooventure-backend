package gemini

import (
	"context"
	"encoding/json"
	"log"
	"os"
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

func (g *GeminiAI) PredictImageAnimal(ctx context.Context, data []byte, typeFile string) model.Animal {
	resp, err := g.model.GenerateContent(ctx,
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
		genai.Text("if the image not animal and not have any animal in the image please just fill the json with not animal"),
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

	var animal model.Animal
	if err := jsoniter.Unmarshal([]byte(cleanedInput), &animal); err != nil {
		log.Fatal(err)
	}

	return animal
}
