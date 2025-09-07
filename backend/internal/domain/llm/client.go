package llm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goda6565/ai-consultant/backend/internal/domain/errors"
)

type Provider string

const (
	OpenAI   Provider = "openai"
	VertexAI Provider = "vertexai"
)

type LLMModel string

const (
	GPT4o         LLMModel = "gpt-4o"
	GPT5          LLMModel = "gpt-5"
	Gemini25Flash LLMModel = "gemini-2.5-flash"
)

type EmbeddingModel string

const (
	EmbeddingModelOpenAIEmbeddings EmbeddingModel = "text-embedding-3-small"
	GeminiEmbedding001             EmbeddingModel = "gemini-embedding-001"
)

type LLMClient interface {
	GenerateText(ctx context.Context, input GenerateTextInput) (*GenerateTextOutput, error)
	GenerateStructuredText(ctx context.Context, input GenerateStructuredTextInput) (*GenerateStructuredTextOutput, error)
	GenerateEmbedding(ctx context.Context, input GenerateEmbeddingInput) (*GenerateEmbeddingOutput, error)
	GenerateEmbeddingBatch(ctx context.Context, input GenerateEmbeddingBatchInput) (*GenerateEmbeddingBatchOutput, error)
}

type LLMConfig struct {
	Provider Provider
	Model    LLMModel
}

func (c LLMConfig) Validate() error {
	switch c.Provider {
	case OpenAI:
		if c.Model != GPT4o && c.Model != GPT5 {
			return errors.NewDomainError(errors.ValidationError, fmt.Sprintf("invalid model %s", c.Model))
		}
	case VertexAI:
		if c.Model != Gemini25Flash {
			return errors.NewDomainError(errors.ValidationError, fmt.Sprintf("invalid model %s", c.Model))
		}
	default:
		return errors.NewDomainError(errors.ValidationError, fmt.Sprintf("invalid provider %s", c.Provider))
	}
	return nil
}

type EmbeddingConfig struct {
	Provider Provider
	Model    EmbeddingModel
}

func (c EmbeddingConfig) Validate() error {
	switch c.Provider {
	case OpenAI:
		if c.Model != EmbeddingModelOpenAIEmbeddings {
			return errors.NewDomainError(errors.ValidationError, fmt.Sprintf("invalid model %s", c.Model))
		}
	case VertexAI:
		if c.Model != GeminiEmbedding001 {
			return errors.NewDomainError(errors.ValidationError, fmt.Sprintf("invalid model %s", c.Model))
		}
	default:
		return errors.NewDomainError(errors.ValidationError, fmt.Sprintf("invalid provider %s", c.Provider))
	}
	return nil
}

type Usage struct {
	InputTokens  int
	OutputTokens int
	TotalTokens  int
}

type GenerateTextInput struct {
	SystemPrompt string
	UserPrompt   string
	Temperature  float32
	Config       LLMConfig
}

type GenerateTextOutput struct {
	Text  string
	Usage Usage
}

type GenerateStructuredTextInput struct {
	SystemPrompt string
	UserPrompt   string
	Temperature  float32
	Schema       json.RawMessage
	Config       LLMConfig
}

type GenerateStructuredTextOutput struct {
	Text  string
	Usage Usage
}

const EmbeddingDimensions = 1536

type GenerateEmbeddingInput struct {
	Text   string
	Config EmbeddingConfig
}

type GenerateEmbeddingOutput struct {
	Embedding []float32
	Usage     Usage
}

type GenerateEmbeddingBatchInput struct {
	Texts  []string
	Config EmbeddingConfig
}

type GenerateEmbeddingBatchOutput struct {
	Embeddings [][]float32
	Usage      Usage
}
