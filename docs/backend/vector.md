## Vector サービス

アプリに登録されたドキュメントを分割・埋め込み計算し、Vector DB（pgvector）に保存する同期サービス。検索は Query 時にベクター類似度で行い、対応するドキュメントのメタ情報を App DB から取得します。

### 役割
- ドキュメントの取り込み・解析（PDF, CSV 等）
- チャンク分割と埋め込み計算
- ベクターテーブルへの保存・削除
- ドキュメント URL 解決（GCS）

### 主なエンドポイント（概略）
- チャンク生成 API: ドキュメント ID を受け取り、OCR/解析→分割→埋め込み→保存まで実施

### 主要コンポーネント
- Router/Handler: `internal/infrastructure/http/echo/vector/*`
- UseCase: `internal/usecase/chunk/*`
- Repository: `internal/infrastructure/google/database/repository/chunk/*`
- SearchClient: `internal/infrastructure/google/database/search/search.go`
- SQLC 定義: `internal/infrastructure/google/database/internal/query/vector/vector.sql`

### データストア
- Vector DB: Postgres + pgvector 拡張（`vectors` テーブル）
- App DB: ドキュメントメタ情報（タイトル、GCS バケット名/オブジェクト名）

### 類似検索
- コサイン類似度（`1 - (embedding <=> $1)`）で降順取得
- 結果に紐づくドキュメント情報を App DB から取得し、`title/content/url` を返却

### 実行方法（ローカル）
- 前提: `.env.vector` に環境変数を設定
- 実行:
```bash
make run-vector
# または
set -a && . .env.vector && set +a && go run main.go vector run
```

### デプロイ
- Cloud Run 上で稼働（README 記載）

### 監視/ログ
- Zap ロガーで処理の各段階を記録
