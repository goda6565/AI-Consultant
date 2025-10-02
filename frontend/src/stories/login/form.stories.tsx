import type { Meta, StoryObj } from "@storybook/nextjs-vite";
import type { z } from "zod";
import type { credentialsSignInSchema } from "@/pages/login/model/zod";
import { LoginForm } from "@/pages/login/ui/form";

const meta: Meta<typeof LoginForm> = {
  title: "Pages/Login/LoginForm",
  component: LoginForm,
};

export default meta;

type Story = StoryObj<typeof LoginForm>;

export const Default: Story = {
  args: {
    onSubmit: async (values: z.infer<typeof credentialsSignInSchema>) => {
      console.log("submitted", values);
    },
  },
};
