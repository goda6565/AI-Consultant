## フロントエンド概要

Next.js 15（App Router、Turbopack）を採用。Storybook による UI 検証、orval による API 型生成を利用しています。認証は Firebase で、API コール時に ID トークンを付与します。

### 実行/ビルド
- 開発: `pnpm dev` （`next dev --turbopack`）
- 本番ビルド: `pnpm build`
- 起動: `pnpm start`
- Storybook: `pnpm storybook`

Docker での実行は `frontend/Dockerfile` を参照。

### ディレクトリ構成（主要）
- `src/pages/*`: 各機能のページ/UI/API hooks
- `src/shared/*`: 共通コンポーネント、API クライアント、env、auth など
- `src/app/*`: ルートレイアウトやプロバイダ
- `src/stories/*`: Storybook 用ストーリー

### API 呼び出し
- クライアントは `src/shared/api/client.ts`
  - `NEXT_PUBLIC_ADMIN_API_URL` / `NEXT_PUBLIC_AGENT_API_URL` を baseURL に設定
  - リクエスト時に Firebase ID Token を `Authorization: Bearer <token>` として付与
- API ラッパ/生成:
  - 管理系: `src/shared/api/admin/*`（orval 生成の型/クライアントを利用）
  - エージェント系: `src/shared/api/agent/*`（存在する場合）

### 認証
- `src/app/provider/auth.tsx` を通じて Firebase Auth コンテキストを提供
- ガード等は `src/shared/auth/*` に配置

### 代表的な画面: 問題詳細/チャット
- UI: `src/pages/problems/ui/page.tsx`（ProblemPage）
  - 取得: `useChatApi(id)` で `problem`/`hearing` を読み込み
  - SSE: `useEventSse({ problemId, enabled })` で進行中イベントを購読
  - メッセージ: `useChatMessage({ hearingId, enabled })` で履歴を管理
  - 実行: `useExecuteHearing(problemId, hearingId)` でメッセージ送信→応答を反映
- 表示: `src/pages/problems/ui/message-view.tsx`（Markdown 対応）
- Storybook: `src/stories/problems/message-view.stories.tsx`

### スタイル/UI
- Tailwind CSS（typography プラグイン）
- `Markdown` コンポーネントで Markdown 描画（`react-markdown` + rehype/remark）

### 型/品質
- TypeScript + Biome（format/lint）
- Steiger による FSD 違反チェック

### 環境変数
- `NEXT_PUBLIC_ADMIN_API_URL`: Admin API のベース URL
- `NEXT_PUBLIC_AGENT_API_URL`: Agent API のベース URL
- Firebase 関連（プロジェクト設定）

### 注意点
- App Router の `use(params)` を利用して動的パラメータを解決
- API エラーは `handleApiError` で正規化しトースト表示
- SSE・Mutation の同時進行で状態競合が起きないよう `mutate`/ローカルメッセージを併用

### アーキテクチャ: Feature-Sliced Design (FSD)
- 本プロジェクトは FSD を採用しています。
- レイヤ順序（依存の矢印は下位→上位）
  - `shared` → `entities` → `features` → `widgets` → `pages` → `app`
- 代表ディレクトリ
  - `src/shared/*`: 基盤（UI、lib、api、config、auth 等）
  - `src/entities/*`: 業務上の最小単位モデル
  - `src/features/*`: 機能ユニット（フォーム、検索など）
  - `src/widgets/*`: 複合 UI セクション
  - `src/pages/*`: ページ単位
  - `src/app/*`: ルート、プロバイダ、レイアウト
- 静的検査
  - `@feature-sliced/steiger-plugin` により境界違反を検出（`pnpm lint`）
  - import は上位レイヤへのみ許可。下位レイヤから `shared` 以外へ逆流しないように実装


