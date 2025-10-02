import type { Meta, StoryObj } from "@storybook/nextjs-vite";
import type { z } from "zod";
import type { MessageFormSchema } from "@/pages/problems/model/zod";
import { MessageForm } from "@/pages/problems/ui/form";

const meta: Meta<typeof MessageForm> = {
  title: "Pages/Problems/MessageForm",
  component: MessageForm,
};

export default meta;

type Story = StoryObj<typeof MessageForm>;

export const Default: Story = {
  args: {
    onSubmit: async (values: z.infer<typeof MessageFormSchema>) => {
      console.log("submitted", values);
    },
    isMutating: false,
  },
};

export const Mutating: Story = {
  args: {
    onSubmit: async (values: z.infer<typeof MessageFormSchema>) => {
      console.log("submitted", values);
    },
    isMutating: true,
  },
};
