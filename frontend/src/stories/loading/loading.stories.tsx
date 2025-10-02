import type { Meta, StoryObj } from "@storybook/nextjs-vite";
import { LoadingPage } from "@/pages/loading/ui/page";

const meta: Meta<typeof LoadingPage> = {
  title: "Pages/Loading/LoadingPage",
  component: LoadingPage,
};

export default meta;

type Story = StoryObj<typeof LoadingPage>;

export const Default: Story = {
  args: {},
};
