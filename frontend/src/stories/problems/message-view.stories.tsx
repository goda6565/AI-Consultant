import type { Meta, StoryObj } from "@storybook/nextjs-vite";
import type { Message } from "@/pages/problems/model/zod";
import { MessageView } from "@/pages/problems/ui/message-view";

const meta: Meta<typeof MessageView> = {
  title: "Pages/Problems/MessageView",
  component: MessageView,
};

export default meta;

type Story = StoryObj<typeof MessageView>;

const sampleMessages: Message[] = [
  {
    role: "assistant",
    message: "こんにちは！どのようなご質問でしょうか？",
  },
  {
    role: "user",
    message: "データベースの設計について相談したいです。",
  },
  {
    role: "assistant",
    message:
      "承知しました。データベースの設計についてお手伝いします。具体的にはどのような内容でしょうか？",
  },
  {
    role: "user",
    message:
      "ECサイトのユーザー管理システムを作っているのですが、ユーザーテーブルの設計で悩んでいます。個人情報と購入履歴をどう分離すべきでしょうか？",
  },
  {
    role: "assistant",
    message:
      "ECサイトの設計ですね。セキュリティとパフォーマンスの観点から、以下のような分離をお勧めします：\n\n1. **usersテーブル**: 基本情報（ID、メールアドレス、パスワードハッシュ等）\n2. **user_profilesテーブル**: 個人情報（氏名、住所、電話番号等）\n3. **ordersテーブル**: 購入履歴\n\nこの構成により、個人情報へのアクセス制御を細かく設定でき、購入履歴の検索パフォーマンスも向上します。",
  },
  {
    role: "user",
    message:
      "なるほど！それぞれのテーブルの関係性はどのように設計すれば良いでしょうか？",
  },
  {
    role: "assistant",
    message:
      "テーブル間の関係性について説明いたします：\n\n**主キーと外部キー**:\n- users.id (主キー)\n- user_profiles.user_id (外部キー) → users.id\n- orders.user_id (外部キー) → users.id\n\n**カーディナリティ**:\n- users : user_profiles = 1:1\n- users : orders = 1:N\n\n**インデックス設計**:\n- user_profiles.user_id にインデックス\n- orders.user_id にインデックス\n- orders.created_at にインデックス（日付範囲検索用）\n\nこれにより効率的なクエリ実行が可能になります。",
  },
  {
    role: "user",
    message:
      "ありがとうございます！セキュリティ面で注意すべき点も教えてください。",
  },
  {
    role: "assistant",
    message:
      "セキュリティ面での重要なポイントをお伝えします：\n\n**1. データ暗号化**\n- 個人情報（氏名、住所、電話番号）は暗号化して保存\n- クレジットカード情報は直接保存せず、トークン化を推奨\n\n**2. アクセス制御**\n- 個人情報テーブルへの直接アクセスを制限\n- アプリケーションレベルでの権限管理\n\n**3. 監査ログ**\n- 個人情報へのアクセス履歴を記録\n- GDPR等の法的要件への対応\n\n**4. データ保持ポリシー**\n- 不要になった個人情報の自動削除機能\n- ユーザーによるデータ削除要求への対応\n\n詳細な実装についても相談いただけますよ。",
  },
  {
    role: "user",
    message:
      "パフォーマンスの最適化についても知りたいです。大量のデータを扱う場合の工夫はありますか？",
  },
  {
    role: "assistant",
    message:
      "大量データを扱う場合の最適化テクニックをご紹介します：\n\n**1. パーティショニング**\n- 購入履歴を日付や地域別に分割\n- 古いデータは別テーブルにアーカイブ\n\n**2. インデックス最適化**\n- 複合インデックスの活用\n- 部分インデックスによるサイズ削減\n\n**3. クエリ最適化**\n- N+1問題の回避\n- 必要な列のみを取得するSELECT\n- ページネーションの実装\n\n**4. キャッシュ戦略**\n- よく使われるユーザー情報のRedisキャッシュ\n- 集計データの事前計算\n\n**5. 読み取り専用レプリカ**\n- 分析用クエリを別DBで実行\n- マスター・スレーブ構成\n\n具体的な実装方法についても詳しく説明できます！",
  },
];

export const Default: Story = {
  args: {
    messages: sampleMessages,
  },
};

export const Empty: Story = {
  args: {
    messages: [],
  },
};
