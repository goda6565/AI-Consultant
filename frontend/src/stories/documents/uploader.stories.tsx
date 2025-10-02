import type { Meta, StoryObj } from "@storybook/nextjs-vite";
import { Uploader } from "@/pages/documents/ui/uploader";
import type { CreateDocumentBody } from "@/shared/api";

const meta: Meta<typeof Uploader> = {
  title: "Pages/Documents/Uploader",
  component: Uploader,
};

export default meta;

type Story = StoryObj<typeof Uploader>;

export const Default: Story = {
  args: {
    trigger: async (data: CreateDocumentBody) => {
      console.log("upload:", data.title, data.documentType);
    },
    isMutating: false,
  },
};

export const Mutating: Story = {
  args: {
    trigger: async (data: CreateDocumentBody) => {
      console.log("upload:", data.title, data.documentType);
    },
    isMutating: true,
  },
};
