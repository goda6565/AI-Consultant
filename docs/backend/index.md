# バックエンド

AI-Consultantプロジェクトのバックエンドシステムです。

## 概要

以下の主要コンポーネントで構成されています：

- **Admin**: 管理機能とAPI提供
- **Agent**: AIエージェント機能
- **Proposal Job**: 提案生成ジョブ処理
- **Vector**: ベクトル検索機能

## 主要機能

- ドキュメント管理とOCR処理
- AI による問題分析と提案生成
- ベクトル検索による関連情報検索
- リアルタイムイベントストリーミング
- 評価システム（LLM-as-a-Judge）

## 技術スタック

- **言語**: Go
- **フレームワーク**: Echo
- **データベース**: PostgreSQL
- **キャッシュ**: Redis (Upstash)
- **クラウド**: Google Cloud Platform
- **コンテナ**: Docker
