## Proposal Job

提案書（Proposal）を自動生成するバッチ/ジョブ。対象の `problemId` を入力として、LLM アクションのオーケストレーションを回し、最終的な提案内容を確定します。

### 役割
- 問題・ヒアリング・フィールド情報を事前取得してエージェント状態を初期化
- Goal 設定 → オーケストレーターで次アクション選択 → アクション実行のループ
- 入出力イベントをイベントストアへ記録し、レポート/アクションを永続化
- 終了時に完了イベントを記録、失敗時は `problem` のステータスを `failed` に更新

### 実行方式
- Cloud Run Jobs を想定（Cloud Run 実行も可）
- コマンド:
```bash
make run-proposal-job
# または
set -a && . .env.proposal-job && set +a && go run main.go proposal-job run
```
- 必須環境変数:
  - `PROBLEM_ID`: 対象の問題 ID

### 主なフロー（概略）
1. 事前取得: Problem/ProblemFields/HearingMessages をロード
2. State 構築: `content`, `history`, `actionHistory` を初期化
3. Goal 生成: LLM によりゴール文生成し State に設定
4. ループ: Orchestrator が `plan/search/analyze/write/review/done` 等から次アクションを決定
5. アクション実行: テンプレート経由で各アクション（検索、要約、執筆等）を実装
6. 永続化: Action/Events を保存、出力があればイベント記録
7. Done 判定: `done` 選択で完了イベントを記録し終了

### 主要コンポーネント
- DI: `di/wire.go` → `InitProposalJob`
- CLI: `cmd/proposal-job.go`
- Job 実装: `internal/infrastructure/job/proposal/*`
- UseCase: `internal/usecase/proposal/*`
- ドメイン（エージェント）: `internal/domain/agent/*`
- 外部連携: Gemini(Vertex AI), Redis(Event), Postgres(App/Vector), Google Search, Vector Search

### エラー処理
- パニック発生時でも Problem ステータスを `failed` に更新
- 外部サービス失敗はインフラ層のエラーラップで原因を記録

### 監視/ログ
- Zap ロガーで各アクション・出力・完了を記録
