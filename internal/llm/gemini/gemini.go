package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"lakpahana/write-my-naita/internal/llm"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiClient struct {
	client *genai.Client
}

func NewGeminiClient(apiKey string) (*GeminiClient, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Error creating Gemini client: %v", err)
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}
	return &GeminiClient{client: client}, nil
}

func (g *GeminiClient) GenerateTimelineForAWeek(prompt string) (llm.WeeklyTrainingReport, error) {

	master := llm.MASTER_PROMPT
	weeklyPrompt := fmt.Sprintf("%s\n%s\nGenerate a weekly training report for the prompt.", master, prompt)

	model := g.client.GenerativeModel("gemini-1.5-flash")

	response, err := model.GenerateContent(context.Background(), genai.Text(
		weeklyPrompt,
	))
	if err != nil {
		return llm.WeeklyTrainingReport{}, fmt.Errorf("failed to generate text: %w", err)
	}

	var report llm.WeeklyTrainingReport
	part := response.Candidates[0].Content.Parts[0]
	partToString := fmt.Sprintf("%v", part)
	log.Printf("partToString: %v", partToString)
	err = json.Unmarshal([]byte(partToString), &report)
	if err != nil {
		return llm.WeeklyTrainingReport{}, fmt.Errorf("failed to parse response into WeeklyTrainingReport: %w", err)
	}
	log.Printf("report: %v", report)
	return report, nil
}
