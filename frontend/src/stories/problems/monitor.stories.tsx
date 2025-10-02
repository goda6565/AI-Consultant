import type { Meta, StoryObj } from "@storybook/nextjs-vite";
import type { Event } from "@/pages/problems/model/zod";
import { Monitor } from "@/pages/problems/ui/monitor";

const meta: Meta<typeof Monitor> = {
  title: "Pages/Problems/Monitor",
  component: Monitor,
};

export default meta;

type Story = StoryObj<typeof Monitor>;

const sampleEvents: Event[] = [
  {
    id: "1",
    eventType: "action",
    actionType: "plan",
    message: "計画を立てています...",
  },
  {
    id: "2",
    eventType: "input",
    actionType: "search",
    message: "関連情報を検索中...",
  },
  {
    id: "3",
    eventType: "output",
    actionType: "analyze",
    message: "分析結果を出力しています",
  },
];

export const Default: Story = {
  args: {
    events: sampleEvents,
    onCopyEvents: () => {
      console.log("Events copied");
    },
  },
};

export const Empty: Story = {
  args: {
    events: [],
    onCopyEvents: () => {
      console.log("Events copied");
    },
  },
};
