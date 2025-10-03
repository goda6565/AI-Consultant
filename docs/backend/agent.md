## Agent サービス

AI Consultant の会話エージェント。ユーザーのヒアリングに応答し、課題解決に向けて対話を進めます。内部では LLM と各種アクションをオーケストレーションして応答を生成します。

### 役割
- ヒアリング実行 API を受け取り、ユーザー発話を元に応答を生成
- 問題 (`problem`) の状態を参照し、処理中・完了状態を適切に制御
- Firebase 認証と OpenAPI バリデーションを通した安全なエンドポイント提供

### 主な入出力
- 入力: `problemId`, `hearingId`, `userMessage`
- 出力: `assistantMessage`, `isCompleted`

概略（内部ハンドラ）:
- `ExecuteHearing` ハンドラが Problem を取得し状態を確認
- `ExecuteHearingUseCase` が LLM・ユースケース群を呼び出し応答文生成
- 結果を返却（必要に応じて対話完了フラグを付与）

### 主要コンポーネント
- Router/Handler: `internal/infrastructure/http/echo/agent/*`
- UseCase: `internal/usecase/hearing/*`
- ドメイン（エージェント）: `internal/domain/agent/*`
- 依存クライアント: Firebase, Gemini (Vertex AI), Postgres

### オーケストレーションの要点
- 目的（Goal）や履歴に基づき、必要なアクション（検索・分析・下書き・レビュー等）を選択
- エージェント用のプロンプト設計と構造化出力で堅牢に制御

### 実行方法（ローカル）
- 前提: `.env.agent` に環境変数を設定
- 実行:
```bash
make run-agent
# または
set -a && . .env.agent && set +a && go run main.go agent run
```

### ランタイム/デプロイ
- Cloud Run 上で稼働（README 記載）
- OpenAPI スキーマに基づくルーティング・バリデーション

### 認証/認可
- Firebase 認証を利用
- 開発環境では Swagger UI（Agent スキーマ）で確認可能（Admin 側と同様の仕組み）

### 代表的なエンドポイント（概略）
- ヒアリング実行（例）
  - 入力: `problemId`, `hearingId`, `body.userMessage`
  - 出力: `assistantMessage`, `isCompleted`
  - 備考: Problem の状態が `processing/done` の場合はエラー応答

### 依存リソース
- App DB (Postgres)
- Vertex AI / Gemini
- Firebase Auth

### ログ/監視
- Zap ロガーを使用
- 重要イベント（アクション選択、完了判定など）を DEBUG/INFO で記録
