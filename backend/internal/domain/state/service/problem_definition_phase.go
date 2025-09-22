package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goda6565/ai-consultant/backend/internal/domain/llm"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
	entity "github.com/goda6565/ai-consultant/backend/internal/domain/state/entity"
	stateValue "github.com/goda6565/ai-consultant/backend/internal/domain/state/value"
	logger "github.com/goda6565/ai-consultant/backend/internal/pkg/logger"
	"github.com/goda6565/ai-consultant/backend/internal/pkg/uuid"
)

type ProblemDefinitionPhase struct {
	llmClient llm.LLMClient
}

func NewProblemDefinitionPhase(llmClient llm.LLMClient) PhaseTemplate {
	return &ProblemDefinitionPhase{llmClient: llmClient}
}

type DefinedProblemStruct struct {
	Summary            string `json:"summary"`
	HearingInformation string `json:"hearingInformation"`
	Goal               string `json:"goal"`
}

type ProblemDefinitionPhaseOutputStruct struct {
	DefinedProblems []DefinedProblemStruct
}

func (p *ProblemDefinitionPhase) Execute(ctx context.Context, state *entity.State) (*entity.State, error) {
	logger := logger.GetLogger(ctx)

	userPrompt := p.createUserPrompt(state)
	systemPrompt := p.createSystemPrompt()
	logger.Info("user prompt", "userPrompt", userPrompt)
	logger.Info("system prompt", "systemPrompt", systemPrompt)

	summary, err := p.llmClient.GenerateStructuredText(ctx, llm.GenerateStructuredTextInput{
		SystemPrompt: systemPrompt,
		UserPrompt:   userPrompt,
		Config:       llm.LLMConfig{Provider: llm.VertexAI, Model: llm.Gemini25Flash},
		Schema: json.RawMessage(`
			{
				"type": "object",
				"properties": {
					"definedProblems": {
						"type": "array",
						"items": {
							"type": "object",
							"properties": {
								"summary": {
									"type": "string"
								},
								"hearingInformation": {
									"type": "string"
								},
								"goal": {
									"type": "string"
								}
							},
							"required": ["summary", "hearingInformation", "goal"]
						}
					}
				},
				"required": ["definedProblems"]
			}
		`),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate summary: %w", err)
	}
	var output ProblemDefinitionPhaseOutputStruct
	if err := json.Unmarshal([]byte(summary.Text), &output); err != nil {
		return nil, fmt.Errorf("failed to unmarshal output: %w", err)
	}
	var definedProblems []entity.ProblemDefinition
	for _, definedProblem := range output.DefinedProblems {
		id, err := sharedValue.NewID(uuid.NewUUID())
		if err != nil {
			return nil, fmt.Errorf("failed to create id: %w", err)
		}
		summary, err := stateValue.NewDefinedProblemSummary(definedProblem.Summary)
		if err != nil {
			return nil, fmt.Errorf("failed to create summary: %w", err)
		}
		hearingInformation, err := stateValue.NewDefinedProblemHearingInformation(definedProblem.HearingInformation)
		if err != nil {
			return nil, fmt.Errorf("failed to create hearing information: %w", err)
		}
		goal, err := stateValue.NewDefinedProblemGoal(definedProblem.Goal)
		if err != nil {
			return nil, fmt.Errorf("failed to create goal: %w", err)
		}
		definedProblems = append(definedProblems, *entity.NewProblemDefinition(id, *summary, *hearingInformation, *goal))
	}
	state.SetProblemDefinitionState(definedProblems)
	return state, nil
}

func (p *ProblemDefinitionPhase) createSystemPrompt() string {
	return `
あなたは熟練した戦略コンサルタントです。
あなたの役割は、クライアントから得たヒアリング情報をもとに、
解決すべき課題（DefinedProblem）を明確に定義することです。

## 出力ルール
- 出力は必ず JSON Schema に従って生成する
- ユーザーが言及した具体的な事実（頻度、数値、条件）は必ず保持する
- 「summary」には課題の要約（1〜2文）
- 「hearingInformation」にはヒアリングで得た主要事実を簡潔にまとめる
- 「goal」には解決後の理想状態を具体的に書く
- 不確実な情報や推測は入れない。会話履歴にない情報は作らない
- 定義する課題は重要なものを優先し、重複は避ける
`
}

func (p *ProblemDefinitionPhase) createUserPrompt(state *entity.State) string {
	return fmt.Sprintf(`
以下はクライアントからのヒアリング結果の整理結果です。

%s

上記をもとに、解決すべき課題を定義してください。
カテゴリごとに重要な課題を抽出し、次の要素を必ず含めてください：

- summary: 1〜2文の課題要約
- hearingInformation: ヒアリングから得た事実（原因や背景）を簡潔に
- goal: この課題が解決された理想状態

出力はJSONのみで、余計な文章や解説は一切書かないでください。
`, state.ToProblemDefinitionPrompt())
}
