## Admin サービス

ドキュメントと問題（Problem）管理を行う REST サービス。OpenAPI ベースのエンドポイントを提供し、開発環境では Swagger UI による検証が可能です。

### 役割
- ドキュメント CRUD（作成/取得/一覧/削除）
- 問題 CRUD、ヒアリング作成/取得、ヒアリングメッセージ一覧
- イベント一覧、レポート取得
- Cloud Tasks による非同期処理トリガ（必要に応じて）

### 主要コンポーネント
- Router/Handler: `internal/infrastructure/http/echo/admin/*`
- UseCase: `internal/usecase/{document,problem,hearing,hearing_message,event,report}/*`
- Repository: `internal/infrastructure/google/database/repository/*`
- 認証: Firebase Auth

### 実行方法（ローカル）
- 前提: `.env.admin` に環境変数を設定
- 実行:
```bash
make run-admin
# または
set -a && . .env.admin && set +a && go run main.go admin run
```

### API スキーマ/Swagger
- OpenAPI スキーマからコード生成（oapi-codegen）
- 開発環境では `/swagger/index.html` で UI を提供

### データストア
- App DB (Postgres): ドキュメント、問題、イベント、レポート等

### 監視/ログ
- Zap ロガーを使用
