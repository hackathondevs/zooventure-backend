package gemini

import (
	"context"
	"fmt"
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
