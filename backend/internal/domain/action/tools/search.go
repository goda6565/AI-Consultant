package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/goda6565/ai-consultant/backend/internal/domain/errors"
	"github.com/goda6565/ai-consultant/backend/internal/domain/llm"
	"github.com/goda6565/ai-consultant/backend/internal/domain/search"
	"github.com/goda6565/ai-consultant/backend/internal/pkg/scraper"

	"github.com/goda6565/ai-consultant/backend/internal/pkg/logger"
)

type SearchTools struct {
	// required
	llmClient llm.LLMClient
	// tools
	WebSearchTool      search.WebSearchClient
	DocumentSearchTool search.DocumentSearchClient
}

type FunctionName string

const (
	FunctionNameWebSearch      FunctionName = "web_search"
	FunctionNameDocumentSearch FunctionName = "document_search"

	WebSearchDescription      = "インターネット全体を対象に情報を検索する。最新情報、業界動向、技術仕様、ベストプラクティス、類似事例、リスク情報、実装手順、最新トレンドなどを調べる場合に使用する。"
	DocumentSearchDescription = "システム内に蓄積されたドキュメントを検索する。社内ナレッジ、過去の事例、内部文書、ユーザー固有の情報、組織固有のベストプラクティスなどを調べる場合に使用する。"
)

const defaultWebSearchMaxNumResults = 5
const defaultDocumentSearchMaxNumResults = 5

func NewSearchTools(llmClient llm.LLMClient, webSearchTool search.WebSearchClient, documentSearchTool search.DocumentSearchClient) *SearchTools {
	return &SearchTools{llmClient: llmClient, WebSearchTool: webSearchTool, DocumentSearchTool: documentSearchTool}
}

func (s *SearchTools) Tools() []llm.Function {
	return []llm.Function{
		{
			Name:        string(FunctionNameWebSearch),
			Description: WebSearchDescription,
			Parameters: json.RawMessage(`{
                "type": "object",
                "properties": {
                    "query": {"type": "string"}
                },
                "required": ["query"]
            }`),
		},
		{
			Name:        string(FunctionNameDocumentSearch),
			Description: DocumentSearchDescription,
			Parameters: json.RawMessage(`{
                "type": "object",
                "properties": {
                    "query": {"type": "string"}
                },
                "required": ["query"]
            }`),
		},
	}
}

type ExecuteInput struct {
	Function llm.FunctionCall
}

type SearchResult struct {
	Title   string
	Content string
	URL     string
}

type ExecuteOutput struct {
	SearchResults []SearchResult
}

func (e *ExecuteOutput) String() string {
	builder := strings.Builder{}
	builder.WriteString("SearchResults:\n")
	for _, result := range e.SearchResults {
		builder.WriteString(fmt.Sprintf("Title: %s\n", result.Title))
		builder.WriteString(fmt.Sprintf("Content: %s\n", result.Content))
		builder.WriteString(fmt.Sprintf("URL: %s\n", result.URL))
	}
	return builder.String()
}

func (s *SearchTools) Execute(ctx context.Context, input ExecuteInput) (*ExecuteOutput, error) {
	functionName := input.Function.Name
	arguments := input.Function.Arguments
	switch functionName {
	case string(FunctionNameWebSearch):
		output, err := s.webSearch(ctx, arguments["query"].(string))
		if err != nil {
			return nil, err
		}
		return &ExecuteOutput{SearchResults: output}, nil
	case string(FunctionNameDocumentSearch):
		output, err := s.documentSearch(ctx, arguments["query"].(string))
		if err != nil {
			return nil, err
		}
		return &ExecuteOutput{SearchResults: output}, nil
	default:
		return nil, errors.NewDomainError(errors.InvalidFunctionName, fmt.Sprintf("invalid function name: %s", functionName))
	}
}

func (s *SearchTools) webSearch(ctx context.Context, query string) ([]SearchResult, error) {
	logger := logger.GetLogger(ctx)
	output, err := s.WebSearchTool.Search(ctx, search.WebSearchInput{Query: query, MaxNumResults: defaultWebSearchMaxNumResults})
	if err != nil {
		return nil, fmt.Errorf("failed to search web: %w", err)
	}

	wg := sync.WaitGroup{}
	scrapeChannel := make(chan SearchResult)
	scrapeClient := scraper.NewScraperClient()
	for _, result := range output.Results {
		wg.Add(1)
		go func(result search.WebSearchResult) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					logger.Error("failed to scrape", "error", r)
				}
			}()
			content, err := scrapeClient.Scrape(ctx, result.URL)
			if err != nil {
				return
			}
			searchResult := SearchResult{Title: result.Title, Content: content, URL: result.URL}
			scrapeChannel <- searchResult
		}(result)
	}

	go func() {
		wg.Wait()
		close(scrapeChannel)
	}()

	searchResults := []SearchResult{}
	for searchResult := range scrapeChannel {
		searchResults = append(searchResults, searchResult)
	}

	return searchResults, nil
}

func (s *SearchTools) documentSearch(ctx context.Context, query string) ([]SearchResult, error) {
	embedding, err := s.llmClient.GenerateEmbedding(ctx, llm.GenerateEmbeddingInput{
		Text: query,
		Config: llm.EmbeddingConfig{
			Provider: llm.VertexAI,
			Model:    llm.GeminiEmbedding001,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding: %w", err)
	}
	output, err := s.DocumentSearchTool.Search(ctx, search.DocumentSearchInput{Query: query, Embedding: &embedding.Embedding, MaxNumResults: defaultDocumentSearchMaxNumResults})
	if err != nil {
		return nil, fmt.Errorf("failed to search document: %w", err)
	}
	searchResults := []SearchResult{}
	for _, result := range output.Results {
		searchResult := SearchResult{Title: result.Title, Content: result.Content, URL: result.URL}
		searchResults = append(searchResults, searchResult)
	}
	return searchResults, nil
}
