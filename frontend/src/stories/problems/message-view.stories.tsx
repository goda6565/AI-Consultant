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
