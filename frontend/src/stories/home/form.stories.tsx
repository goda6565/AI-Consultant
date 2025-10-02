import type { Meta, StoryObj } from "@storybook/nextjs-vite";
import type { z } from "zod";
import type { problemFormSchema } from "@/pages/home/model/zod";
import { ProblemForm } from "@/pages/home/ui/form";

const meta: Meta<typeof ProblemForm> = {
  title: "Pages/Home/ProblemForm",
  component: ProblemForm,
};

export default meta;

type Story = StoryObj<typeof ProblemForm>;

export const Default: Story = {
  args: {
    onSubmit: async (values: z.infer<typeof problemFormSchema>) => {
      console.log("submitted", values);
    },
    isMutating: false,
  },
};
