package service

import (
	"context"
	"fmt"

	"github.com/goda6565/ai-consultant/backend/internal/domain/action/tools"
	actionValue "github.com/goda6565/ai-consultant/backend/internal/domain/action/value"
	agentState "github.com/goda6565/ai-consultant/backend/internal/domain/agent/state"
	"github.com/goda6565/ai-consultant/backend/internal/domain/llm"
)

type SearchAction struct {
	llmClient   llm.LLMClient
	searchTools *tools.SearchTools
}

func NewSearchAction(llmClient llm.LLMClient, searchTools *tools.SearchTools) SearchActionInterface {
	return &SearchAction{llmClient: llmClient, searchTools: searchTools}
}

func (s *SearchAction) Execute(ctx context.Context, input ActionTemplateInput) (*ActionTemplateOutput, error) {
	systemPrompt := s.createSystemPrompt()
	userPrompt := s.createUserPrompt(input.State)
	llmInput := llm.GenerateFunctionCallInput{
		SystemPrompt: systemPrompt,
		UserPrompt:   userPrompt,
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
	action, err := CreateAction(input.State, actionValue.ActionTypeSearch, llmOutput.FunctionCallString(), searchResults.String())
	if err != nil {
		return nil, fmt.Errorf("failed to create action: %w", err)
	}
	return &ActionTemplateOutput{Action: *action, Content: input.State.GetContent()}, nil // search action does not change content
}

func (s *SearchAction) createSystemPrompt() string {
	return searchSystemPrompt
}

func (s *SearchAction) createUserPrompt(state agentState.State) string {
	return fmt.Sprintf(searchUserPrompt, state.ToPrompt())
}

var searchSystemPrompt = `
あなたは「問題解決エージェント」の検索担当です。  
目的は、現在の課題解決に不足している情報を補うため、適切な検索を実行することです。  

ルール:
- 出力は必ず関数呼び出しのみとする。文章を直接生成してはいけない。  
- 利用可能な検索方法は渡されたリストの中から選ぶ。  
- どの情報が不足しているかを判断し、その不足を埋める検索クエリを設計する。  
- クエリは具体的かつ明確にする（例: "クラウド移行 セキュリティ ベストプラクティス 2025"）。  
- URLや既存リンクを改変してはいけない。  
- 不確実な情報は推測せず、「調査が必要」と判断したら必ず検索を行う。  
`

var searchUserPrompt = `
以下は現在のエージェント状態です。  
この情報を踏まえて、不足している部分を補うための検索を必ず実行してください。  

- 不足がなければ「検索不要」の関数呼び出しを返してください。
- 出力は必ず関数呼び出し形式のみで返してください。

=== 現在の状態 ===
%s
`
