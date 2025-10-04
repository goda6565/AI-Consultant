package proposaljob

import (
	"time"

	hearingEntity "github.com/goda6565/ai-consultant/backend/internal/domain/hearing/entity"
	hearingMessageEntity "github.com/goda6565/ai-consultant/backend/internal/domain/hearing_message/entity"
	hearingMessageValue "github.com/goda6565/ai-consultant/backend/internal/domain/hearing_message/value"
	problemEntity "github.com/goda6565/ai-consultant/backend/internal/domain/problem/entity"
	problemValue "github.com/goda6565/ai-consultant/backend/internal/domain/problem/value"
	problemFieldEntity "github.com/goda6565/ai-consultant/backend/internal/domain/problem_field/entity"
	problemFieldValue "github.com/goda6565/ai-consultant/backend/internal/domain/problem_field/value"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
	"github.com/google/uuid"
)

type MockProblemData struct {
	Title       string
	Description string
	Status      problemValue.Status
}

type MockProblemField struct {
	Text     string
	Answered bool
}

type MockHearingMessage struct {
	FieldIndex int
	Role       string
	Message    string
}

var mockProblemData = MockProblemData{
	Title:       "銀行支店CS改善実行計画",
	Description: "大手銀行の特定支店において「顧客満足度（CS）」が低下している。原因を特定し、短期・中期・長期の改善策とKPIを伴う実行計画を作成せよ。",
	Status:      problemValue.StatusProcessing,
}

var mockProblemFields = []MockProblemField{
	{
		Text:     "顧客満足度調査結果詳細（設問別/属性別/時系列）",
		Answered: true,
	},
	{
		Text:     "顧客クレーム・問い合わせ内容と件数（時系列/分類別）",
		Answered: true,
	},
	{
		Text:     "支店業務プロセスと顧客接点（窓口/ATM/相談ブース）",
		Answered: true,
	},
	{
		Text:     "従業員の業務負荷とスキルレベル（部署別/役割別）",
		Answered: true,
	},
	{
		Text:     "競合他社のCS水準と提供サービス",
		Answered: true,
	},
}

var mockHearingMessages = []MockHearingMessage{
	// Field 0: 顧客満足度調査結果詳細
	{
		FieldIndex: 0,
		Role:       "assistant",
		Message:    "現在、貴行の特定支店において、顧客満足度に関する調査はどのような形で実施されていますでしょうか？また、その調査結果は、設問別、顧客属性別（例：年齢層、取引内容など）、時系列で分析できる形で整理されていますか？",
	},
	{
		FieldIndex: 0,
		Role:       "user",
		Message:    "当支店では、半年に一度、全顧客を対象としたアンケート調査を実施しています。調査は郵送とインターネットの両方で回答いただける形式で、窓口対応、待ち時間、手続きの分かりやすさ、デジタルサービスの利便性などを設問にしています。結果は本部システムで自動集計され、設問別・属性別・時系列での比較が可能です。",
	},

	// Field 1: クレーム・問い合わせ
	{
		FieldIndex: 1,
		Role:       "assistant",
		Message:    "貴行の特定支店では、顧客からのクレームや問い合わせは、どのようなチャネルで受け付けていらっしゃいますか？また、それらの内容はどのように記録・管理されていますか？",
	},
	{
		FieldIndex: 1,
		Role:       "user",
		Message:    "当支店では、窓口、電話、メール、ウェブサイト経由で受け付けています。窓口や電話の内容は専用システムに入力し、案件ごとに管理番号を付与。『待ち時間』『窓口対応』『手続き関連』『デジタルサービス』などに分類して月次集計しています。",
	},

	// Field 2: 支店業務プロセスと顧客接点
	{
		FieldIndex: 2,
		Role:       "assistant",
		Message:    "顧客が支店を訪れる際、どのような目的で、どのような流れでサービスを受けていることが多いでしょうか？",
	},
	{
		FieldIndex: 2,
		Role:       "user",
		Message:    "来店目的は①預金や振込②住所変更などの各種手続き③口座開設・ローン相談に大別されます。受付で目的確認→番号札→待合→窓口対応の流れです。入出金で1〜5分、住所変更で10〜20分、ローン相談で30〜60分かかります。",
	},
	{
		FieldIndex: 2,
		Role:       "assistant",
		Message:    "お客様が待ち時間を経験するのはどのような状況が多いでしょうか？また、その際にどのような配慮をされていますか？",
	},
	{
		FieldIndex: 2,
		Role:       "user",
		Message:    "月曜・金曜午前や昼前後（11〜14時）に集中し、平均15〜20分待ち、時に30分超。受付で待ち時間目安を提示し、デジタルサイネージやスマホ通知で案内しています。",
	},

	// Field 3: 従業員の業務負荷とスキルレベル
	{
		FieldIndex: 3,
		Role:       "assistant",
		Message:    "窓口や相談業務を担当される方々の業務負荷はどの程度でしょうか？また、負荷の種類に違いはありますか？",
	},
	{
		FieldIndex: 3,
		Role:       "user",
		Message:    "窓口担当は1日30〜40件（1件1〜20分）、残業月10〜15時間。相談担当は3〜5件（30分〜1時間）で心理的・知識的負荷が高い傾向です。",
	},
	{
		FieldIndex: 3,
		Role:       "assistant",
		Message:    "顧客対応の質や効率に影響を与えているスキル面の課題はありますか？",
	},
	{
		FieldIndex: 3,
		Role:       "user",
		Message:    "窓口担当は複雑手続きの正確性とデジタルツール活用力、相談担当は商品知識の習熟度と提案力が課題です。",
	},

	// Field 4: 競合他社のCS水準と提供サービス
	{
		FieldIndex: 4,
		Role:       "assistant",
		Message:    "顧客満足度の観点で、特に意識されている競合他社はございますか？",
	},
	{
		FieldIndex: 4,
		Role:       "user",
		Message:    "地域密着型の地方銀行（〇〇銀行、△△信用金庫）やネット銀行（楽天銀行、住信SBIなど）を意識しています。地方銀は待ち時間の短さ、ネット銀は利便性に強みがあります。メガバンク（三菱UFJ、三井住友）も競合です。",
	},
	{
		FieldIndex: 4,
		Role:       "assistant",
		Message:    "ヒアリングが完了しました。",
	},
}

// MockDataProvider provides mock data for evaluation
type MockDataProvider struct {
	now               time.Time
	problemFieldIDs   []string
	hearingID         string
	hearingMessageIDs []string
}

// NewMockDataProvider creates a new mock data provider
func NewMockDataProvider() *MockDataProvider {
	// Generate UUIDs for all entities
	var problemFieldIDs []string
	for range mockProblemFields {
		problemFieldIDs = append(problemFieldIDs, uuid.New().String())
	}

	var hearingMessageIDs []string
	for range mockHearingMessages {
		hearingMessageIDs = append(hearingMessageIDs, uuid.New().String())
	}

	return &MockDataProvider{
		now:               time.Now(),
		problemFieldIDs:   problemFieldIDs,
		hearingID:         uuid.New().String(),
		hearingMessageIDs: hearingMessageIDs,
	}
}

// CreateMockProblem creates a mock problem entity
func (m *MockDataProvider) CreateMockProblem() *problemEntity.Problem {
	id, _ := sharedValue.NewID(EvaluateProblemID)
	title, _ := problemValue.NewTitle(mockProblemData.Title)
	description, _ := problemValue.NewDescription(mockProblemData.Description)
	status := problemValue.StatusProcessing

	return problemEntity.NewProblem(id, *title, *description, status, &m.now)
}

// CreateMockProblemFields creates mock problem field entities
func (m *MockDataProvider) CreateMockProblemFields() []problemFieldEntity.ProblemField {
	problemID, _ := sharedValue.NewID(EvaluateProblemID)

	var fields []problemFieldEntity.ProblemField

	for i, fieldData := range mockProblemFields {
		fieldID, _ := sharedValue.NewID(m.problemFieldIDs[i])
		field, _ := problemFieldValue.NewField(fieldData.Text)
		answered := problemFieldValue.NewAnswered(fieldData.Answered)

		fields = append(fields, *problemFieldEntity.NewProblemField(fieldID, problemID, *field, *answered, &m.now))
	}

	return fields
}

// CreateMockHearing creates a mock hearing entity
func (m *MockDataProvider) CreateMockHearing() *hearingEntity.Hearing {
	hearingID, _ := sharedValue.NewID(m.hearingID)
	problemID, _ := sharedValue.NewID(EvaluateProblemID)

	return hearingEntity.NewHearing(hearingID, problemID, &m.now)
}

// CreateMockHearingMessages creates mock hearing message entities
func (m *MockDataProvider) CreateMockHearingMessages() []hearingMessageEntity.HearingMessage {
	hearingID, _ := sharedValue.NewID(m.hearingID)

	var messages []hearingMessageEntity.HearingMessage

	for i, msgData := range mockHearingMessages {
		msgID, _ := sharedValue.NewID(m.hearingMessageIDs[i])
		fieldIndex := msgData.FieldIndex
		fieldID, _ := sharedValue.NewID(m.problemFieldIDs[fieldIndex])
		message, _ := hearingMessageValue.NewMessage(msgData.Message)
		role, _ := hearingMessageValue.NewRole(msgData.Role)

		messages = append(messages, *hearingMessageEntity.NewHearingMessage(msgID, hearingID, fieldID, role, *message, &m.now))
	}

	return messages
}

// GetMockData returns all mock data needed for evaluation
func (m *MockDataProvider) GetMockData(problemId string) (
	*problemEntity.Problem,
	[]problemFieldEntity.ProblemField,
	*hearingEntity.Hearing,
	[]hearingMessageEntity.HearingMessage,
) {
	return m.CreateMockProblem(),
		m.CreateMockProblemFields(),
		m.CreateMockHearing(),
		m.CreateMockHearingMessages()
}
