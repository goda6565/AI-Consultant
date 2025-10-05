"use client";

import { DocumentTable } from "@/pages/documents/ui/document-table";
import { Uploader } from "@/pages/documents/ui/uploader";
import {
  deleteDocument,
  useCreateDocument,
  useListDocuments,
} from "@/shared/api";
import { Heading } from "@/shared/ui";

export function DocumentPage() {
  const { data, isLoading, mutate } = useListDocuments();
  const { trigger, isMutating } = useCreateDocument();
  return (
    <div className="flex flex-col gap-5">
      <div className="flex gap-2 justify-between items-center">
        <Heading>Documents</Heading>
        <Uploader trigger={trigger} isMutating={isMutating} />
      </div>
      <DocumentTable
        documents={data?.documents || []}
        isLoading={isLoading}
        mutate={mutate}
        deleteDocument={deleteDocument}
      />
    </div>
  );
}
