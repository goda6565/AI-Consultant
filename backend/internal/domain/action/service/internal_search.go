package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/goda6565/ai-consultant/backend/internal/domain/action/tools"
	actionValue "github.com/goda6565/ai-consultant/backend/internal/domain/action/value"
	agentState "github.com/goda6565/ai-consultant/backend/internal/domain/agent/state"
	"github.com/goda6565/ai-consultant/backend/internal/domain/llm"
	"github.com/goda6565/ai-consultant/backend/internal/domain/prompt/service"
	"github.com/goda6565/ai-consultant/backend/internal/pkg/logger"
)

const maxInternalSearchDecomposeTopics = 5

type InternalSearchAction struct {
	llmClient     llm.LLMClient
	searchTools   *tools.SearchTools
	promptBuilder *service.PromptBuilder
}

func NewInternalSearchAction(llmClient llm.LLMClient, searchTools *tools.SearchTools, promptBuilder *service.PromptBuilder) InternalSearchActionInterface {
	return &InternalSearchAction{llmClient: llmClient, searchTools: searchTools, promptBuilder: promptBuilder}
}

func (s *InternalSearchAction) Execute(ctx context.Context, input ActionTemplateInput) (*ActionTemplateOutput, error) {
	logger := logger.GetLogger(ctx)
	// 1. decompose
	topics, err := s.decompose(ctx, InternalSearchDecomposeInput{
		MaxTopics: maxInternalSearchDecomposeTopics,
		State:     input.State,
	})
	logger.Debug("decompose", "topics", topics)
	if err != nil {
		return nil, fmt.Errorf("failed to decompose: %w", err)
	}

	// 2. explore
	wg := sync.WaitGroup{}
	results := []string{}
	resultChannel := make(chan string, len(topics.SearchTopics))

	for _, topic := range topics.SearchTopics {
		wg.Add(1)
		go func(topic string) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					logger.Error("failed to explore", "error", r)
				}
			}()
			result, err := s.explore(ctx, InternalSearchExploreInput{
				Topic: topic,
			})
			if err != nil {
				logger.Error("failed to explore", "error", err)
			}
			resultChannel <- result.result
		}(topic)
	}

	go func() {
		wg.Wait()
		close(resultChannel)
	}()

	for result := range resultChannel {
		results = append(results, result)
	}

	logger.Debug("explore", "results", results)

	// 3. synthesize
	synthesizedResults := []string{}
	for _, result := range results {
		synthesizedResult, err := s.synthesize(ctx, InternalSearchSynthesizeInput{
			Result: result,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to synthesize: %w", err)
		}
		logger.Debug("synthesize", "result", synthesizedResult.result)
		synthesizedResults = append(synthesizedResults, synthesizedResult.result)
	}

	action, err := CreateAction(input.State, actionValue.ActionTypeInternalSearch, "", strings.Join(synthesizedResults, "\n"))
	if err != nil {
		return nil, fmt.Errorf("failed to create action: %w", err)
	}

	return &ActionTemplateOutput{Action: *action, Content: input.State.GetContent()}, nil // search action does not change content
}

type InternalSearchDecomposeInput struct {
	MaxTopics int
	State     agentState.State
}

type InternalSearchDecomposeOutput struct {
	SearchTopics []string `json:"searchTopics"`
}

func (s *InternalSearchAction) decompose(ctx context.Context, input InternalSearchDecomposeInput) (*InternalSearchDecomposeOutput, error) {
	prompt := s.promptBuilder.Build(service.PromptBuilderInput{
		Name:       "decompose",
		ActionType: actionValue.ActionTypeInternalSearch,
		State:      input.State,
		Input:      fmt.Sprintf("%d", input.MaxTopics),
	})
	llmInput := llm.GenerateStructuredTextInput{
		SystemPrompt: prompt.SystemPrompt,
		UserPrompt:   prompt.UserPrompt,
		Config:       llm.LLMConfig{Provider: llm.VertexAI, Model: llm.Gemini25Flash},
		Temperature:  0.0,
		Schema: json.RawMessage(`
			{
				"type": "object",
				"properties": {
					"searchTopics": {
						"type": "array",
						"items": {
							"type": "string"
						}
					}
				},
				"required": ["searchTopics"]
			}
		`),
	}
	llmOutput, err := s.llmClient.GenerateStructuredText(ctx, llmInput)
	if err != nil {
		return nil, fmt.Errorf("failed to generate text: %w", err)
	}
	var output InternalSearchDecomposeOutput
	if err := json.Unmarshal([]byte(llmOutput.Text), &output); err != nil {
		return nil, fmt.Errorf("failed to unmarshal output: %w", err)
	}
	return &output, nil
}

type InternalSearchExploreInput struct {
	Topic string
}

type InternalSearchExploreOutput struct {
	result string
}

func (s *InternalSearchAction) explore(ctx context.Context, input InternalSearchExploreInput) (*InternalSearchExploreOutput, error) {
	prompt := s.promptBuilder.Build(service.PromptBuilderInput{
		Name:       "explore",
		ActionType: actionValue.ActionTypeInternalSearch,
		Input:      input.Topic,
	})
	llmInput := llm.GenerateFunctionCallInput{
		SystemPrompt: prompt.SystemPrompt,
		UserPrompt:   prompt.UserPrompt,
		Config:       llm.LLMConfig{Provider: llm.VertexAI, Model: llm.Gemini25Flash},
		Temperature:  0.0,
		Functions:    s.searchTools.Tools(),
	}
	llmOutput, err := s.llmClient.GenerateFunctionCall(ctx, llmInput)
	if err != nil {
		return nil, fmt.Errorf("failed to generate function call: %w", err)
	}
	searchResults, err := s.searchTools.Execute(ctx, tools.ExecuteInput{Function: llmOutput.FunctionCall})
	if err != nil {
		return nil, fmt.Errorf("failed to execute search tools: %w", err)
	}
	return &InternalSearchExploreOutput{result: searchResults.String()}, nil
}

type InternalSearchSynthesizeInput struct {
	Result string
}

type InternalSearchSynthesizeOutput struct {
	result string
}

func (s *InternalSearchAction) synthesize(ctx context.Context, input InternalSearchSynthesizeInput) (*InternalSearchSynthesizeOutput, error) {
	prompt := s.promptBuilder.Build(service.PromptBuilderInput{
		Name:       "synthesize",
		ActionType: actionValue.ActionTypeInternalSearch,
		Input:      input.Result,
	})
	llmInput := llm.GenerateTextInput{
		SystemPrompt: prompt.SystemPrompt,
		UserPrompt:   prompt.UserPrompt,
		Config:       llm.LLMConfig{Provider: llm.VertexAI, Model: llm.Gemini25Flash},
		Temperature:  0.0,
	}
	llmOutput, err := s.llmClient.GenerateText(ctx, llmInput)
	if err != nil {
		return nil, fmt.Errorf("failed to generate text: %w", err)
	}
	return &InternalSearchSynthesizeOutput{result: llmOutput.Text}, nil
}
