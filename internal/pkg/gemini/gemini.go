package gemini

import (
	"context"
	"encoding/json"
	"fmt"
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

func (g *GeminiAI) PredictImageAnimal(ctx context.Context, picture []byte) (string, error) {
	content, err := g.model.GenerateContent(ctx,
		genai.Text("What type of animal in this picture you think it is?"),
		genai.Text("Give me the exact singular answer without anything else, if its not animal just send me \"notanimal\""),
		genai.Text("Answer it in indonesian language"),
		genai.ImageData("jpeg", picture),
	)
	if err != nil {
		return "", err
	}
	var formattedContent strings.Builder
	if content != nil && content.Candidates != nil {
		for _, cand := range content.Candidates {
			if cand.Content != nil {
				for _, part := range cand.Content.Parts {
					formattedContent.WriteString(fmt.Sprintf("%v", part))
				}
			}
		}
	}
	return formattedContent.String(), nil
}

func (g *GeminiAI) GenerateTrivia(ctx context.Context, animal string) (model.Trivia, error) {
	content, err := g.model.GenerateContent(ctx,
		genai.Text("Make me a single trivia quiz about "+animal+" with 3 possible answer, use the following JSON structure for the output"),
		genai.Text(`{"question"; "Whats the question?","answer":"the correct answer value","optA"; "Answer Option A","optB"; "Answer Option B","optC"; "Answer Option C"}`),
		genai.Text("Give me only the JSON output in one-line, without anything else"),
		genai.Text("Answer it in indonesian language"),
	)
	if err != nil {
		return model.Trivia{}, err
	}
	part := content.Candidates[0].Content.Parts[0]
	jsonByte, err := json.Marshal(part)
	if err != nil {
		log.Fatal(err)
	}
	jsonStr, err := strconv.Unquote(string(jsonByte))
	if err != nil {
		return model.Trivia{}, err
	}
	var trivia model.Trivia
	if err := jsoniter.Unmarshal([]byte(jsonStr), &trivia); err != nil {
		return model.Trivia{}, err
	}
	return trivia, nil
}
