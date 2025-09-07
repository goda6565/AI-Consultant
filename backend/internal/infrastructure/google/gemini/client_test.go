package gemini

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/goda6565/ai-consultant/backend/internal/domain/llm"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/environment"
)

type Response struct {
	Capital string `json:"capital" required:"true"`
}

func TestGeminiClient_GenerateStructuredText(t *testing.T) {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT_ID")
	if projectID == "" {
		t.Skip("GOOGLE_CLOUD_PROJECT_ID is not set")
	}
	location := os.Getenv("VERTEX_AI_LOCATION")
	if location == "" {
		t.Skip("VERTEX_AI_LOCATION is not set")
	}
	client := NewGeminiClient(context.Background(), &environment.Environment{
		GoogleCloudEnvironment: environment.GoogleCloudEnvironment{
			ProjectID: projectID,
		},
		VertexAIEnvironment: environment.VertexAIEnvironment{
			VertexAILocation: location,
		},
	})
	response, err := client.GenerateStructuredText(context.Background(), llm.GenerateStructuredTextInput{
		SystemPrompt: "You are a helpful assistant.",
		UserPrompt:   "What is the capital of France?",
		Schema: json.RawMessage(`
			{
				"type": "object",
				"properties": {
					"capital": {
						"type": "string"
					}
				},
				"required": ["capital"]
			}
		`),
		Temperature: 0.0,
		Config: llm.LLMConfig{
			Provider: llm.VertexAI,
			Model:    llm.Gemini25Flash,
		},
	})
	if err != nil {
		t.Fatalf("failed to generate structured text: %v", err)
	}
	t.Logf("response: %v", response.Text)
	var parsed Response
	if err := json.Unmarshal([]byte(response.Text), &parsed); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	t.Logf("response: %v", parsed.Capital)
}
