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

func (s *SearchAction) createSystemPrompt() string {
	return searchSystemPrompt
}

func (s *SearchAction) createUserPrompt(state agentState.State) string {
	return fmt.Sprintf(searchUserPrompt, state.ToPrompt())
}

var searchSystemPrompt = `
あなたは「問題解決エージェント」の検索担当です。  
目的は、現在の課題解決に不足している情報を補うため、適切な検索を実行することです。  
ユーザーが提供する「現在の状態」には、最終ゴールが明示されています。
この最終ゴールを常に参照しながら、検索を行ってください。

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
