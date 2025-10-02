import type { Meta, StoryObj } from "@storybook/nextjs-vite";
import { DocumentTable } from "@/pages/documents/ui/document-table";
import type { Document as DocumentResponse } from "@/shared/api";
import { DocumentStatus, DocumentType } from "@/shared/api/admin/model";

const meta: Meta<typeof DocumentTable> = {
  title: "Pages/Documents/DocumentTable",
  component: DocumentTable,
};

export default meta;

type Story = StoryObj<typeof DocumentTable>;

const now = new Date().toISOString();

const sampleDocuments: DocumentResponse[] = [
  {
    id: "doc-1",
    title: "要件定義",
    documentType: DocumentType.markdown,
    bucketName: "bucket-a",
    objectName: "a/req.md",
    documentStatus: DocumentStatus.processing,
    retryCount: 0,
    createdAt: now,
    updatedAt: now,
  },
  {
    id: "doc-2",
    title: "設計書",
    documentType: DocumentType.pdf,
    bucketName: "bucket-b",
    objectName: "b/design.pdf",
    documentStatus: DocumentStatus.done,
    retryCount: 0,
    createdAt: now,
    updatedAt: now,
  },
  {
    id: "doc-3",
    title: "データ定義",
    documentType: DocumentType.csv,
    bucketName: "bucket-c",
    objectName: "c/schema.csv",
    documentStatus: DocumentStatus.failed,
    retryCount: 1,
    createdAt: now,
    updatedAt: now,
  },
  {
    id: "doc-4",
    title: "データ定義",
    documentType: DocumentType.csv,
    bucketName: "bucket-c",
    objectName: "c/schema.csv",
    documentStatus: DocumentStatus.failed,
    retryCount: 1,
    createdAt: now,
    updatedAt: now,
  },
  {
    id: "doc-5",
    title: "データ定義",
    documentType: DocumentType.csv,
    bucketName: "bucket-c",
    objectName: "c/schema.csv",
    documentStatus: DocumentStatus.failed,
    retryCount: 1,
    createdAt: now,
    updatedAt: now,
  },
  {
    id: "doc-6",
    title: "データ定義",
    documentType: DocumentType.csv,
    bucketName: "bucket-c",
    objectName: "c/schema.csv",
    documentStatus: DocumentStatus.failed,
    retryCount: 1,
    createdAt: now,
    updatedAt: now,
  },
  {
    id: "doc-7",
    title: "データ定義",
    documentType: DocumentType.csv,
    bucketName: "bucket-c",
    objectName: "c/schema.pdf",
    documentStatus: DocumentStatus.processing,
    retryCount: 1,
    createdAt: now,
    updatedAt: now,
  },
  {
    id: "doc-8",
    title: "データ定義",
    documentType: DocumentType.csv,
    bucketName: "bucket-c",
    objectName: "c/schema.csv",
    documentStatus: DocumentStatus.processing,
    retryCount: 1,
    createdAt: now,
    updatedAt: now,
  },
  {
    id: "doc-9",
    title: "データ定義",
    documentType: DocumentType.csv,
    bucketName: "bucket-c",
    objectName: "c/schema.pdf",
    documentStatus: DocumentStatus.processing,
    retryCount: 1,
    createdAt: now,
    updatedAt: now,
  },
  {
    id: "doc-10",
    title: "データ定義",
    documentType: DocumentType.csv,
    bucketName: "bucket-c",
    objectName: "c/schema.pdf",
    documentStatus: DocumentStatus.done,
    retryCount: 1,
    createdAt: now,
    updatedAt: now,
  },
];

export const Default: Story = {
  args: {
    documents: sampleDocuments,
    isLoading: false,
    mutate: async () => {
      console.log("mutate");
    },
    deleteDocument: async (id: string) => {
      console.log("delete:", id);
      return null;
    },
  },
};

export const Loading: Story = {
  args: {
    documents: [],
    isLoading: true,
    mutate: async () => {
      console.log("mutate");
    },
    deleteDocument: async (_id: string) => {
      return null;
    },
  },
};

export const Empty: Story = {
  args: {
    documents: [],
    isLoading: false,
    mutate: async () => {
      console.log("mutate");
    },
    deleteDocument: async (_id: string) => {
      return null;
    },
  },
};
