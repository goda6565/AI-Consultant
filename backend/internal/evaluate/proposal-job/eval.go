package proposaljob

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	actionEntity "github.com/goda6565/ai-consultant/backend/internal/domain/action/entity"
	actionRepository "github.com/goda6565/ai-consultant/backend/internal/domain/action/repository"
	actionService "github.com/goda6565/ai-consultant/backend/internal/domain/action/service"
	agentService "github.com/goda6565/ai-consultant/backend/internal/domain/agent/service"
	mockEventRepository "github.com/goda6565/ai-consultant/backend/internal/domain/event/repository/mock"
	mockHearingRepository "github.com/goda6565/ai-consultant/backend/internal/domain/hearing/repository/mock"
	mockHearingMessageRepository "github.com/goda6565/ai-consultant/backend/internal/domain/hearing_message/repository/mock"
	mockJobConfigRepository "github.com/goda6565/ai-consultant/backend/internal/domain/job_config/repository/mock"
	mockProblemRepository "github.com/goda6565/ai-consultant/backend/internal/domain/problem/repository/mock"
	mockProblemFieldRepository "github.com/goda6565/ai-consultant/backend/internal/domain/problem_field/repository/mock"
	reportEntity "github.com/goda6565/ai-consultant/backend/internal/domain/report/entity"
	reportRepository "github.com/goda6565/ai-consultant/backend/internal/domain/report/repository"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
	"github.com/goda6565/ai-consultant/backend/internal/evaluate"
	llmasjudge "github.com/goda6565/ai-consultant/backend/internal/evaluate/proposal-job/llm-as-a-judge"
	"github.com/goda6565/ai-consultant/backend/internal/pkg/logger"
	"github.com/goda6565/ai-consultant/backend/internal/usecase/proposal"
	"go.uber.org/mock/gomock"
)

const EvaluateProblemID = "378b8576-1482-4d14-852f-90e753fabce3"
const numEvaluations = 3

type ProposalJobEval struct {
	orchestrator     *agentService.Orchestrator
	summarizeService *agentService.SummarizeService
	goalService      *agentService.GoalService
	terminator       *agentService.Terminator
	skipper          *agentService.Skipper
	actionFactory    *actionService.ActionFactory
	reportRepository reportRepository.ReportRepository
	actionRepository actionRepository.ActionRepository
	judge            *llmasjudge.Judge
	outputDir        string
}

func NewProposalJobEval(
	orchestrator *agentService.Orchestrator,
	summarizeService *agentService.SummarizeService,
	goalService *agentService.GoalService,
	terminator *agentService.Terminator,
	skipper *agentService.Skipper,
	actionFactory *actionService.ActionFactory,
	reportRepository reportRepository.ReportRepository,
	actionRepository actionRepository.ActionRepository,
	judge *llmasjudge.Judge,
) evaluate.Evaluator {
	return &ProposalJobEval{
		orchestrator:     orchestrator,
		summarizeService: summarizeService,
		goalService:      goalService,
		terminator:       terminator,
		skipper:          skipper,
		actionFactory:    actionFactory,
		reportRepository: reportRepository,
		actionRepository: actionRepository,
		judge:            judge,
		outputDir:        "",
	}
}

func (e *ProposalJobEval) prepareForEvaluation(ctx context.Context) (proposal.ExecuteProposalInputPort, error) {
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	// モックデータプロバイダーを作成
	mockProvider := NewMockDataProvider()
	mockProblem, mockProblemFields, mockHearing, mockHearingMessages, mockJobConfig := mockProvider.GetMockData(EvaluateProblemID)

	problemRepository := mockProblemRepository.NewMockProblemRepository(ctrl)
	problemRepository.EXPECT().FindById(gomock.Any(), gomock.Any()).Return(mockProblem, nil).AnyTimes()
	problemRepository.EXPECT().UpdateStatus(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	problemFieldRepository := mockProblemFieldRepository.NewMockProblemFieldRepository(ctrl)
	problemFieldRepository.EXPECT().FindByProblemID(gomock.Any(), gomock.Any()).Return(mockProblemFields, nil).AnyTimes()

	hearingRepository := mockHearingRepository.NewMockHearingRepository(ctrl)
	hearingRepository.EXPECT().FindByProblemId(gomock.Any(), gomock.Any()).Return(mockHearing, nil).AnyTimes()

	hearingMessageRepository := mockHearingMessageRepository.NewMockHearingMessageRepository(ctrl)
	hearingMessageRepository.EXPECT().FindByHearingID(gomock.Any(), gomock.Any()).Return(mockHearingMessages, nil).AnyTimes()

	eventRepository := mockEventRepository.NewMockEventRepository(ctrl)
	eventRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	jobConfigRepository := mockJobConfigRepository.NewMockJobConfigRepository(ctrl)
	jobConfigRepository.EXPECT().FindByProblemID(gomock.Any(), gomock.Any()).Return(mockJobConfig, nil).AnyTimes()

	executeProposalUseCase := proposal.NewExecuteProposalUseCase(
		problemRepository,
		problemFieldRepository,
		hearingRepository,
		hearingMessageRepository,
		e.actionRepository,
		eventRepository,
		e.orchestrator,
		e.summarizeService,
		e.goalService,
		e.terminator,
		e.skipper,
		e.actionFactory,
		e.reportRepository,
		jobConfigRepository,
	)
	return executeProposalUseCase, nil
}

// Evaluate executes the proposal evaluation
func (e *ProposalJobEval) Execute(ctx context.Context) error {
	logger := logger.GetLogger(ctx)
	executeProposalUseCase, err := e.prepareForEvaluation(ctx)
	if err != nil {
		logger.Error("failed to prepare for evaluation", "error", err)
		return err
	}

	input := proposal.ExecuteProposalUseCaseInput{
		ProblemID: EvaluateProblemID,
	}

	startTime := time.Now()

	err = executeProposalUseCase.Execute(ctx, input)
	if err != nil {
		logger.Error("failed to execute proposal", "error", err)
		return err
	}

	duration := time.Since(startTime)

	problemID, err := sharedValue.NewID(EvaluateProblemID)
	if err != nil {
		logger.Error("failed to find report", "error", err)
		return err
	}

	report, err := e.reportRepository.FindByProblemID(ctx, problemID)
	if err != nil {
		return err
	}

	if report == nil {
		logger.Error("report not found")
		return errors.New("report not found")
	}

	actions, err := e.actionRepository.FindByProblemID(ctx, problemID)
	if err != nil {
		logger.Error("failed to find actions", "error", err)
		return err
	}

	if len(actions) == 0 {
		logger.Error("actions not found")
		return errors.New("actions not found")
	}

	err = e.outputResults(ctx, report, actions, duration)
	if err != nil {
		logger.Error("failed to output results", "error", err)
		return fmt.Errorf("failed to output results: %w", err)
	}

	if err := e.runLLMJudgment(ctx, report, actions); err != nil {
		logger.Error("failed to run LLM-as-a-Judge evaluation", "error", err)
	}

	return nil
}

// outputResults creates a timestamped directory and outputs actions.csv and report.md
func (e *ProposalJobEval) outputResults(ctx context.Context, report *reportEntity.Report, actions []actionEntity.Action, duration time.Duration) error {
	timestamp := time.Now().Format("20060102_150405")

	baseDir := filepath.Join("internal", "evaluate", "proposal-job", "outputs")
	outputDir := filepath.Join(baseDir, fmt.Sprintf("evaluation_results_%s", timestamp))

	// 出力ディレクトリを保存（LLM-as-a-Judgeで再利用）
	e.outputDir = outputDir

	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	err = e.outputActionsCSV(filepath.Join(outputDir, "actions.csv"), actions)
	if err != nil {
		return fmt.Errorf("failed to output actions.csv: %w", err)
	}

	err = e.outputReportMD(filepath.Join(outputDir, "report.md"), report)
	if err != nil {
		return fmt.Errorf("failed to output report.md: %w", err)
	}

	err = e.outputMetadata(filepath.Join(outputDir, "metadata.txt"), duration, len(actions))
	if err != nil {
		return fmt.Errorf("failed to output metadata: %w", err)
	}

	fmt.Printf("評価結果を出力しました: %s (実行時間: %s, アクション数: %d)\n", outputDir, duration.Round(time.Second), len(actions))
	return nil
}

// outputActionsCSV outputs actions to CSV file
func (e *ProposalJobEval) outputActionsCSV(filename string, actions []actionEntity.Action) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			_ = closeErr // ignore close error as file operations are already completed
		}
	}()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{"ID", "ProblemID", "ActionType", "Input", "Output", "CreatedAt"}); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	for _, action := range actions {
		createdAt := ""
		if action.GetCreatedAt() != nil {
			createdAt = action.GetCreatedAt().Format(time.RFC3339)
		}

		record := []string{
			action.GetID().Value(),
			action.GetProblemID().Value(),
			action.GetActionType().Value(),
			func() string {
				input := action.GetInput()
				return input.Value()
			}(),
			func() string {
				output := action.GetOutput()
				return output.Value()
			}(),
			createdAt,
		}

		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write action to CSV: %w", err)
		}
	}

	return nil
}

// outputReportMD outputs report to Markdown file
func (e *ProposalJobEval) outputReportMD(filename string, report *reportEntity.Report) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create markdown file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			_ = closeErr // ignore close error as file operations are already completed
		}
	}()

	content := report.GetContent()
	_, err = file.WriteString(content.Value())
	if err != nil {
		return fmt.Errorf("failed to write to markdown file: %w", err)
	}

	return nil
}

// runLLMJudgment runs LLM-as-a-Judge evaluation
func (e *ProposalJobEval) runLLMJudgment(ctx context.Context, report *reportEntity.Report, actions []actionEntity.Action) error {
	fmt.Println("\n🤖 LLM-as-a-Judge評価を開始します...")

	mockProvider := NewMockDataProvider()
	mockProblem, _, _, _, _ := mockProvider.GetMockData(EvaluateProblemID)
	problemDescription := mockProblem.GetDescription().Value()

	startTime := time.Now()
	result, err := e.judge.Judge(ctx, problemDescription, actions, report, numEvaluations)
	if err != nil {
		return fmt.Errorf("failed to judge: %w", err)
	}
	judgmentDuration := time.Since(startTime)

	if err := e.outputJudgmentResult(result, judgmentDuration); err != nil {
		return fmt.Errorf("failed to output judgment result: %w", err)
	}

	return nil
}

// outputJudgmentResult outputs the judgment result to the same directory as other results
func (e *ProposalJobEval) outputJudgmentResult(result *llmasjudge.JudgmentResult, duration time.Duration) error {
	// 既存の出力ディレクトリを使用（outputResultsで作成されたディレクトリ）
	if e.outputDir == "" {
		return fmt.Errorf("output directory not set")
	}

	// JSON形式で詳細結果を出力
	if err := e.outputJudgmentJSON(filepath.Join(e.outputDir, "llm_judgment.json"), result, duration); err != nil {
		return fmt.Errorf("failed to output JSON: %w", err)
	}

	// 人間が読みやすいテキスト形式でも出力
	if err := e.outputJudgmentText(filepath.Join(e.outputDir, "llm_judgment.txt"), result, duration); err != nil {
		return fmt.Errorf("failed to output text: %w", err)
	}

	// コンソールにサマリーを表示
	e.printJudgmentSummary(result, e.outputDir, duration)

	return nil
}

// outputJudgmentJSON outputs judgment result as JSON
func (e *ProposalJobEval) outputJudgmentJSON(filename string, result *llmasjudge.JudgmentResult, duration time.Duration) error {
	output := map[string]interface{}{
		"result":          result,
		"evaluation_time": duration.String(),
		"timestamp":       time.Now().Format(time.RFC3339),
	}

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	return nil
}

// outputJudgmentText outputs judgment result as text
func (e *ProposalJobEval) outputJudgmentText(filename string, result *llmasjudge.JudgmentResult, duration time.Duration) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create text file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			_ = closeErr // ignore close error as file operations are already completed
		}
	}()

	if _, err := fmt.Fprintf(file, "LLM-as-a-Judge 評価結果\n"); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}
	if _, err := fmt.Fprintf(file, "======================\n\n"); err != nil {
		return fmt.Errorf("failed to write separator: %w", err)
	}
	if _, err := fmt.Fprintf(file, "評価日時: %s\n", time.Now().Format("2006-01-02 15:04:05")); err != nil {
		return fmt.Errorf("failed to write timestamp: %w", err)
	}
	if _, err := fmt.Fprintf(file, "評価時間: %s\n\n", duration.Round(time.Millisecond)); err != nil {
		return fmt.Errorf("failed to write duration: %w", err)
	}

	if _, err := fmt.Fprintf(file, "## スコア\n\n"); err != nil {
		return fmt.Errorf("failed to write score header: %w", err)
	}
	if _, err := fmt.Fprintf(file, "総合評価:           %d/10\n", result.OverallScore); err != nil {
		return fmt.Errorf("failed to write overall score: %w", err)
	}
	if _, err := fmt.Fprintf(file, "信頼性・参考文献:   %d/10\n", result.ReliabilityCheck); err != nil {
		return fmt.Errorf("failed to write reliability score: %w", err)
	}
	if _, err := fmt.Fprintf(file, "論理的整合性:       %d/10\n", result.LogicalConsistency); err != nil {
		return fmt.Errorf("failed to write logical consistency score: %w", err)
	}
	if _, err := fmt.Fprintf(file, "実用性・実行可能性: %d/10\n", result.Practicality); err != nil {
		return fmt.Errorf("failed to write practicality score: %w", err)
	}
	if _, err := fmt.Fprintf(file, "網羅性・完成度:     %d/10\n", result.Completeness); err != nil {
		return fmt.Errorf("failed to write completeness score: %w", err)
	}
	if _, err := fmt.Fprintf(file, "アクションの適切さ: %d/10\n", result.ActionAppropriateness); err != nil {
		return fmt.Errorf("failed to write action appropriateness score: %w", err)
	}
	if _, err := fmt.Fprintf(file, "目標達成度:         %d/10\n\n", result.GoalAchievement); err != nil {
		return fmt.Errorf("failed to write goal achievement score: %w", err)
	}

	if _, err := fmt.Fprintf(file, "## 目標達成\n\n"); err != nil {
		return fmt.Errorf("failed to write goal achievement header: %w", err)
	}
	if result.GoalAchieved {
		if _, err := fmt.Fprintf(file, "✅ 達成\n\n"); err != nil {
			return fmt.Errorf("failed to write goal achieved: %w", err)
		}
	} else {
		if _, err := fmt.Fprintf(file, "❌ 未達成\n\n"); err != nil {
			return fmt.Errorf("failed to write goal not achieved: %w", err)
		}
	}

	if _, err := fmt.Fprintf(file, "## 総合評価\n\n%s\n\n", result.Summary); err != nil {
		return fmt.Errorf("failed to write summary: %w", err)
	}
	if _, err := fmt.Fprintf(file, "## 良かった点\n\n%s\n\n", result.Strengths); err != nil {
		return fmt.Errorf("failed to write strengths: %w", err)
	}
	if _, err := fmt.Fprintf(file, "## 改善点\n\n%s\n\n", result.Weaknesses); err != nil {
		return fmt.Errorf("failed to write weaknesses: %w", err)
	}
	if _, err := fmt.Fprintf(file, "## 判定理由\n\n%s\n", result.Reasoning); err != nil {
		return fmt.Errorf("failed to write reasoning: %w", err)
	}

	return nil
}

// printJudgmentSummary prints judgment summary to console
func (e *ProposalJobEval) printJudgmentSummary(result *llmasjudge.JudgmentResult, outputDir string, duration time.Duration) {
	goalSymbol := "❌"
	if result.GoalAchieved {
		goalSymbol = "✅"
	}

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🤖 LLM-as-a-Judge 評価結果")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("総合評価: %d/10 | 目標達成: %s\n", result.OverallScore, goalSymbol)
	fmt.Println(strings.Repeat("-", 60))
	fmt.Printf("信頼性・参考文献:   %d/10\n", result.ReliabilityCheck)
	fmt.Printf("論理的整合性:       %d/10\n", result.LogicalConsistency)
	fmt.Printf("実用性・実行可能性: %d/10\n", result.Practicality)
	fmt.Printf("網羅性・完成度:     %d/10\n", result.Completeness)
	fmt.Printf("アクションの適切さ: %d/10\n", result.ActionAppropriateness)
	fmt.Printf("目標達成度:         %d/10\n", result.GoalAchievement)
	fmt.Println(strings.Repeat("-", 60))
	fmt.Printf("評価時間: %s\n", duration.Round(time.Millisecond))
	fmt.Printf("結果保存先: %s\n", outputDir)
	fmt.Println(strings.Repeat("=", 60) + "\n")
}

// outputMetadata outputs evaluation metadata
func (e *ProposalJobEval) outputMetadata(filename string, duration time.Duration, actionCount int) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create metadata file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			_ = closeErr // ignore close error as file operations are already completed
		}
	}()

	metadata := "Evaluation Metadata\n"
	metadata += "===================\n\n"
	metadata += fmt.Sprintf("Problem ID: %s\n", EvaluateProblemID)
	metadata += fmt.Sprintf("Execution Time: %s\n", duration.Round(time.Millisecond))
	metadata += fmt.Sprintf("Action Count: %d\n", actionCount)
	metadata += fmt.Sprintf("Timestamp: %s\n", time.Now().Format(time.RFC3339))

	_, err = file.WriteString(metadata)
	if err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}

	return nil
}
