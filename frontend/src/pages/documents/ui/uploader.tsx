"use client";

import { useRef } from "react";
import { toast } from "sonner";
import {
  type CreateDocumentBody,
  DocumentType,
  type ErrorResponse,
} from "@/shared/api";
import { Button } from "@/shared/ui";

type AvailableDocumentType = {
  extension: string;
  name: string;
};

const AVAILABLE_DOCUMENT_TYPES: AvailableDocumentType[] = [
  { extension: "pdf", name: "PDF" },
  { extension: "md", name: "Markdown" },
  { extension: "csv", name: "CSV" },
];

const MAX_FILE_SIZE = {
  size: 10 * 1024 * 1024,
  name: "10MB",
};

const MAX_PDF_PAGE_COUNT = 15;

const getFileExtension = (file: File): string => {
  return file.name.split(".").pop() || "";
};

const validateFile = async (file: File): Promise<boolean> => {
  console.log(file);
  if (
    !AVAILABLE_DOCUMENT_TYPES.some(
      (type) => type.extension === getFileExtension(file),
    )
  ) {
    toast.error(
      `無効なファイルです利用できるファイルは${AVAILABLE_DOCUMENT_TYPES.map((type) => type.name).join(", ")}です`,
    );
    return false;
  }
  if (file.size > MAX_FILE_SIZE.size) {
    toast.error(
      `ファイルサイズが大きすぎます最大サイズは${MAX_FILE_SIZE.name}です`,
    );
    return false;
  }
  if (file.type === "application/pdf") {
    const pageCount = await getPdfPageCount(file);
    if (pageCount > MAX_PDF_PAGE_COUNT) {
      toast.error(
        `PDFは最大${MAX_PDF_PAGE_COUNT}ページまでです（${pageCount}ページ）`,
      );
      return false;
    }
  }
  return true;
};

const toBase64 = (file: File): Promise<string> => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = (event) => {
      const result = event.target?.result as string;
      const base64 = result.split(",")[1];
      resolve(base64);
    };
    reader.onerror = reject;
    reader.readAsDataURL(file);
  });
};

const getPdfPageCount = async (file: File): Promise<number> => {
  const buffer = await file.arrayBuffer();
  const text = new TextDecoder("latin1").decode(buffer);
  const matches = text.match(/\/Type\s*\/Page\b/g);
  return matches ? matches.length : 0;
};

const getDocumentType = (file: File): DocumentType => {
  switch (getFileExtension(file)) {
    case "pdf":
      return DocumentType.pdf;
    case "md":
      return DocumentType.markdown;
    case "csv":
      return DocumentType.csv;
  }
  throw new Error("Invalid file type");
};

type UploaderProps = {
  trigger: (data: CreateDocumentBody) => void;
  isMutating: boolean;
};

export function Uploader({ trigger, isMutating }: UploaderProps) {
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleButtonClick = () => {
    fileInputRef.current?.click();
  };

  return (
    <div className="max-w-2xs">
      <input
        ref={fileInputRef}
        type="file"
        className="hidden"
        onChange={async (e) => {
          const file = e.target.files?.[0];
          if (file && (await validateFile(file))) {
            try {
              const base64 = await toBase64(file);
              const documentType = getDocumentType(file);
              const title = file.name.replace(/\.[^/.]+$/, "");
              await trigger({
                title,
                documentType: documentType,
                data: base64,
              });
              toast.success("ファイルをアップロードしました");
            } catch (error) {
              if (error && typeof error === "object" && "message" in error) {
                const errorResponse = error as ErrorResponse;
                toast.error(
                  `アップロードに失敗しました: ${errorResponse.message}`,
                );
              } else {
                toast.error("ファイルのアップロードに失敗しました");
              }
            }
          }
          e.target.value = "";
        }}
      />
      <Button
        onClick={handleButtonClick}
        variant="outline"
        disabled={isMutating}
      >
        {isMutating ? "アップロード中..." : "ファイルを追加"}
      </Button>
    </div>
  );
}
