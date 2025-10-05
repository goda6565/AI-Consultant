"use client";

import type { ColumnDef } from "@tanstack/react-table";
import type { Document } from "./document";

export const columns: ColumnDef<Document>[] = [
  {
    accessorKey: "title",
    header: "Title",
  },
  {
    accessorKey: "documentType",
    header: "Document Type",
  },
  {
    accessorKey: "documentStatus",
    header: "Document Status",
  },
  {
    accessorKey: "createdAt",
    header: "Created At",
  },
];
