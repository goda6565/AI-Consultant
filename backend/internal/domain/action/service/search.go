package service

import (
	"context"
	"fmt"

	"github.com/goda6565/ai-consultant/backend/internal/domain/action/tools"
	actionValue "github.com/goda6565/ai-consultant/backend/internal/domain/action/value"
	"github.com/goda6565/ai-consultant/backend/internal/domain/llm"
	"github.com/goda6565/ai-consultant/backend/internal/domain/prompt/service"
)

type SearchAction struct {
	llmClient     llm.LLMClient
	searchTools   *tools.SearchTools
	promptBuilder *service.PromptBuilder
}

func NewSearchAction(llmClient llm.LLMClient, searchTools *tools.SearchTools, promptBuilder *service.PromptBuilder) SearchActionInterface {
	return &SearchAction{llmClient: llmClient, searchTools: searchTools, promptBuilder: promptBuilder}
}

func (s *SearchAction) Execute(ctx context.Context, input ActionTemplateInput) (*ActionTemplateOutput, error) {
	prompt := s.promptBuilder.Build(service.PromptBuilderInput{
		ActionType: actionValue.ActionTypeSearch,
		State:      input.State,
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
	summarizedSearchResults, err := s.summarizeSearchAction(ctx, llmOutput.FunctionCallValueString(), *searchResults)
	if err != nil {
		return nil, fmt.Errorf("failed to summarize search results: %w", err)
	}
	action, err := CreateAction(input.State, actionValue.ActionTypeSearch, llmOutput.FunctionCallString(), summarizedSearchResults)
	if err != nil {
		return nil, fmt.Errorf("failed to create action: %w", err)
	}
	return &ActionTemplateOutput{Action: *action, Content: input.State.GetContent()}, nil // search action does not change content
}

func (s *SearchAction) summarizeSearchAction(ctx context.Context, params string, searchResults tools.ExecuteOutput) (string, error) {
	systemPrompt := s.createSummarizeSearchActionSystemPrompt()
	userPrompt := s.createSummarizeSearchActionUserPrompt(params, searchResults)
	llmInput := llm.GenerateTextInput{
		SystemPrompt: systemPrompt,
		UserPrompt:   userPrompt,
		Config:       llm.LLMConfig{Provider: llm.VertexAI, Model: llm.Gemini25Flash},
		Temperature:  0.0,
	}
	llmOutput, err := s.llmClient.GenerateText(ctx, llmInput)
	if err != nil {
		return "", fmt.Errorf("failed to generate text: %w", err)
	}
	return llmOutput.Text, nil
}

func (s *SearchAction) createSummarizeSearchActionSystemPrompt() string {
	return summarizeSearchActionSystemPrompt
}

func (s *SearchAction) createSummarizeSearchActionUserPrompt(params string, searchResults tools.ExecuteOutput) string {
	return fmt.Sprintf(summarizeSearchActionUserPrompt, params, searchResults.String())
}

var summarizeSearchActionSystemPrompt = `
あなたは「問題解決エージェント」の検索結果要約担当です。  
目的は、検索ツールから得られた情報を整理し、課題解決に役立つ要点だけを簡潔にまとめることです。  

ルール:
- 検索結果の信頼できる情報を中心にまとめる。  
- 参照URLは削除・改変せず、必ず出力に残すこと。  
- URLは対応する情報の直後に記載する。  
- 重複する内容は統合して要点のみ残す。  
- 不明確・推測的な部分はそのまま「情報不足」と明示する。  
- 出力はMarkdown形式とし、見出し・箇条書きを用いて整理する。  
- 新しい情報を勝手に追加したり、検索結果を改変してはいけない。  
- まとめは、現在の課題解決に役立つ観点（例: メリット、リスク、具体的手順など）を優先する。  
 - 検索結果が現在の課題と無関係、または有用な情報が見つからない場合は、出力全体を「関連する情報はなかった」の一文のみとする（他の文章や装飾は不要）。  
`

var summarizeSearchActionUserPrompt = `
以下は検索ツールから得られた結果です。  
これらを基に、課題解決に役立つ要点をMarkdown形式で整理してください。  

ルール:
- 各要点には必ず対応するURLを残すこと（削除・改変禁止）。  
- 情報の要点を簡潔にまとめる。  
- 重複は統合する。  
- 不足している部分は「情報不足」と明示する。  
- 検索クエリに関連する情報がない場合は、出力全体を「関連する情報はなかった」の一文のみとする（他の文章や装飾は不要）。  

=== 検索クエリ ===
%s

=== 検索結果 ===
%s
`
