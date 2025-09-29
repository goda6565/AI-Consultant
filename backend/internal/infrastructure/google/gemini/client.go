package gemini

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goda6565/ai-consultant/backend/internal/domain/llm"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/environment"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/errors"
	"google.golang.org/genai"
)

type GeminiClient struct {
	client *genai.Client
}

func NewGeminiClient(ctx context.Context, e *environment.Environment) llm.LLMClient {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		Backend: genai.Backend(
			genai.BackendVertexAI,
		),
		Project:  e.ProjectID,
		Location: e.VertexAILocation,
	})
	if err != nil {
		panic(err)
	}
	return &GeminiClient{client: client}
}

func (c *GeminiClient) GenerateText(ctx context.Context, input llm.GenerateTextInput) (*llm.GenerateTextOutput, error) {
	response, err := c.client.Models.GenerateContent(ctx, string(input.Config.Model), []*genai.Content{
		genai.NewContentFromText(input.SystemPrompt, genai.RoleModel),
		genai.NewContentFromText(input.UserPrompt, genai.RoleUser),
	}, &genai.GenerateContentConfig{
		Temperature: &input.Temperature,
	})
	if err != nil {
		return nil, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to generate text: %v", err))
	}
	usage := llm.Usage{
		InputTokens:  int(response.UsageMetadata.PromptTokenCount),
		OutputTokens: int(response.UsageMetadata.CandidatesTokenCount),
		TotalTokens:  int(response.UsageMetadata.TotalTokenCount),
	}
	return &llm.GenerateTextOutput{Text: response.Text(), Usage: usage}, nil
}

func (c *GeminiClient) GenerateStructuredText(ctx context.Context, input llm.GenerateStructuredTextInput) (*llm.GenerateStructuredTextOutput, error) {
	var schema genai.Schema
	if err := json.Unmarshal(input.Schema, &schema); err != nil {
		return nil, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to unmarshal schema: %v", err))
	}

	config := &genai.GenerateContentConfig{
		Temperature:      &input.Temperature,
		ResponseMIMEType: "application/json",
		ResponseSchema:   &schema,
	}

	response, err := c.client.Models.GenerateContent(ctx, string(input.Config.Model), []*genai.Content{
		genai.NewContentFromText(input.SystemPrompt, genai.RoleModel),
		genai.NewContentFromText(input.UserPrompt, genai.RoleUser),
	}, config)
	if err != nil {
		return nil, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to generate structured text: %v", err))
	}
	usage := llm.Usage{
		InputTokens:  int(response.UsageMetadata.PromptTokenCount),
		OutputTokens: int(response.UsageMetadata.CandidatesTokenCount),
		TotalTokens:  int(response.UsageMetadata.TotalTokenCount),
	}
	return &llm.GenerateStructuredTextOutput{Text: response.Text(), Usage: usage}, nil
}

func (c *GeminiClient) GenerateFunctionCall(ctx context.Context, input llm.GenerateFunctionCallInput) (*llm.GenerateFunctionCallOutput, error) {
	// Build tools with function declarations
	var fns []*genai.FunctionDeclaration
	for _, fn := range input.Functions {
		var schema genai.Schema
		if err := json.Unmarshal(fn.Parameters, &schema); err != nil {
			return nil, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to unmarshal function schema for %s: %v", fn.Name, err))
		}
		fns = append(fns, &genai.FunctionDeclaration{
			Name:        fn.Name,
			Description: fn.Description,
			Parameters:  &schema,
		})
	}

	tools := []*genai.Tool{}
	if len(fns) > 0 {
		tools = append(tools, &genai.Tool{FunctionDeclarations: fns})
	}

	cfg := &genai.GenerateContentConfig{
		Temperature: &input.Temperature,
		Tools:       tools,
		ToolConfig: &genai.ToolConfig{
			FunctionCallingConfig: &genai.FunctionCallingConfig{
				// ANY mode forces the model to predict only function calls
				Mode: genai.FunctionCallingConfigModeAny,
			},
		},
	}

	response, err := c.client.Models.GenerateContent(ctx, string(input.Config.Model), []*genai.Content{
		genai.NewContentFromText(input.SystemPrompt, genai.RoleModel),
		genai.NewContentFromText(input.UserPrompt, genai.RoleUser),
	}, cfg)
	if err != nil {
		return nil, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to generate function call: %v", err))
	}

	// Extract function call if present
	if len(response.Candidates) == 0 || response.Candidates[0] == nil || response.Candidates[0].Content == nil {
		return nil, errors.NewInfrastructureError(errors.ExternalServiceError, "no function call found")
	}

	parts := response.Candidates[0].Content.Parts
	if len(parts) == 0 || parts[0] == nil || parts[0].FunctionCall == nil {
		return nil, errors.NewInfrastructureError(errors.ExternalServiceError, "no function call found")
	}

	functionCall := llm.FunctionCall{
		Name:      parts[0].FunctionCall.Name,
		Arguments: parts[0].FunctionCall.Args,
	}

	usage := llm.Usage{
		InputTokens:  int(response.UsageMetadata.PromptTokenCount),
		OutputTokens: int(response.UsageMetadata.CandidatesTokenCount),
		TotalTokens:  int(response.UsageMetadata.TotalTokenCount),
	}

	return &llm.GenerateFunctionCallOutput{FunctionCall: functionCall, Usage: usage}, nil
}

func (c *GeminiClient) GenerateEmbedding(ctx context.Context, input llm.GenerateEmbeddingInput) (*llm.GenerateEmbeddingOutput, error) {
	contents := []*genai.Content{
		genai.NewContentFromText(input.Text, genai.RoleUser),
	}
	dimensions := int32(llm.EmbeddingDimensions)
	response, err := c.client.Models.EmbedContent(ctx, string(input.Config.Model), contents, &genai.EmbedContentConfig{OutputDimensionality: &dimensions})
	if err != nil {
		return nil, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to generate embedding: %v", err))
	}
	embeddings := response.Embeddings[0].Values
	usage := llm.Usage{
		InputTokens:  0,
		OutputTokens: 0,
		TotalTokens:  int(response.Metadata.BillableCharacterCount),
	}
	return &llm.GenerateEmbeddingOutput{Embedding: embeddings, Usage: usage}, nil
}

func (c *GeminiClient) GenerateEmbeddingBatch(ctx context.Context, input llm.GenerateEmbeddingBatchInput) (*llm.GenerateEmbeddingBatchOutput, error) {
	contents := make([]*genai.Content, len(input.Texts))
	for i, text := range input.Texts {
		contents[i] = genai.NewContentFromText(text, genai.RoleUser)
	}
	dimensions := int32(llm.EmbeddingDimensions)
	response, err := c.client.Models.EmbedContent(ctx, string(input.Config.Model), contents, &genai.EmbedContentConfig{OutputDimensionality: &dimensions})
	if err != nil {
		return nil, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to generate embedding batch: %v", err))
	}
	embeddings := make([][]float32, len(response.Embeddings))
	for i, embedding := range response.Embeddings {
		embeddings[i] = embedding.Values
	}
	usage := llm.Usage{
		InputTokens:  0,
		OutputTokens: 0,
		TotalTokens:  int(response.Metadata.BillableCharacterCount),
	}
	return &llm.GenerateEmbeddingBatchOutput{Embeddings: embeddings, Usage: usage}, nil
}

func (c *GeminiClient) GetTokenCount(ctx context.Context, input llm.CountTokenInput) (*llm.CountTokenOutput, error) {
	contents := []*genai.Content{
		genai.NewContentFromText(input.Text, genai.RoleUser),
	}
	response, err := c.client.Models.CountTokens(ctx, string(input.Config.Model), contents, &genai.CountTokensConfig{})
	if err != nil {
		return nil, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to get token count: %v", err))
	}
	return &llm.CountTokenOutput{TokenCount: int(response.TotalTokens)}, nil
}
