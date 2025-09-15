"use client";

import { DocumentTable } from "@/pages/documents/ui/document-table";
import { Uploader } from "@/pages/documents/ui/uploader";
import { Heading } from "@/shared/ui";

export function DocumentPage() {
  return (
    <div className="flex flex-col gap-5">
      <div className="flex gap-2 justify-between items-center">
        <Heading>Documents</Heading>
        <Uploader />
      </div>
      <DocumentTable />
    </div>
  );
}
