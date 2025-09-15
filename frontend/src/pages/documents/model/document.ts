import type { Document as DocumentResponse } from "@/shared/api";

export type Document = {
  id: string;
  title: string;
  documentType: "pdf" | "markdown" | "csv";
  documentStatus: "pending" | "processing" | "done" | "failed";
  createdAt: string;
};

const formatCreatedAt = (isoString: string): string => {
  const date = new Date(isoString);
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, "0");
  const day = String(date.getDate()).padStart(2, "0");
  const hour = String(date.getHours()).padStart(2, "0");
  const minute = String(date.getMinutes()).padStart(2, "0");

  return `${year}年${month}月${day}日 ${hour}時${minute}分`;
};

export const responseToDocument = (response: DocumentResponse): Document => {
  return {
    id: response.id,
    title: response.title,
    documentType: response.documentType as Document["documentType"],
    documentStatus: response.documentStatus as Document["documentStatus"],
    createdAt: formatCreatedAt(response.createdAt),
  };
};
