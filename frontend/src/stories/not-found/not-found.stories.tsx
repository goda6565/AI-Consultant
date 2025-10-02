import type { Meta, StoryObj } from "@storybook/nextjs-vite";
import { NotFound } from "@/pages/not-found/ui/page";

const meta: Meta<typeof NotFound> = {
  title: "Pages/NotFound/NotFound",
  component: NotFound,
};

export default meta;

type Story = StoryObj<typeof NotFound>;

export const Default: Story = {
  args: {},
};
