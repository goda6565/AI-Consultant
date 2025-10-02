import type { Meta, StoryObj } from "@storybook/nextjs-vite";
import { ReportView } from "@/pages/problems/ui/report-view";

const meta: Meta<typeof ReportView> = {
  title: "Pages/Problems/ReportView",
  component: ReportView,
};

export default meta;

type Story = StoryObj<typeof ReportView>;

const sampleReport = {
  id: "1",
  problemId: "problem-1",
  content: `# レポート

## 概要
このレポートは問題の分析結果をまとめたものです。

## 詳細
- ポイント1: データベース設計の改善
- ポイント2: パフォーマンスの最適化
- ポイント3: セキュリティ強化

## 結論
提案された改善案を実装することで、システムの品質が向上します。`,
  createdAt: "2025-10-02T12:00:00Z",
};

export const Default: Story = {
  args: {
    report: sampleReport,
    isLoading: false,
    onCopyReport: () => {
      console.log("Report copied");
    },
  },
};

export const Loading: Story = {
  args: {
    report: undefined,
    isLoading: true,
    onCopyReport: () => {
      console.log("Report copied");
    },
  },
};

export const NoReport: Story = {
  args: {
    report: undefined,
    isLoading: false,
    onCopyReport: () => {
      console.log("Report copied");
    },
  },
};
