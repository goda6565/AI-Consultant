"use client";

import {
  type ColumnDef,
  flexRender,
  getCoreRowModel,
  useReactTable,
} from "@tanstack/react-table";
import { useEffect, useMemo, useState } from "react";
import { toast } from "sonner";
import { columns } from "@/pages/documents/model/columns";
import type { Document } from "@/pages/documents/model/document";
import { responseToDocument } from "@/pages/documents/model/document";
import type { Document as DocumentResponse } from "@/shared/api";
import {
  Button,
  Input,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/shared/ui";

const columnsWithSelection: ColumnDef<Document, unknown>[] = [
  {
    id: "select",
    size: 32,
    header: ({ table }) => (
      <input
        type="checkbox"
        aria-label="Select all"
        checked={table.getIsAllRowsSelected()}
        onChange={table.getToggleAllRowsSelectedHandler()}
      />
    ),
    cell: ({ row }) => (
      <input
        type="checkbox"
        aria-label="Select row"
        checked={row.getIsSelected()}
        onChange={row.getToggleSelectedHandler()}
      />
    ),
    enableSorting: false,
    enableHiding: false,
  },
  ...columns,
];

type DocumentTableProps = {
  documents: DocumentResponse[];
  isLoading: boolean;
  mutate: () => void;
  deleteDocument: (id: string) => void;
};

export function DocumentTable({
  documents,
  isLoading,
  mutate,
  deleteDocument,
}: DocumentTableProps) {
  const [query, setQuery] = useState("");
  const [status, setStatus] = useState<string>("");
  const [docType, setDocType] = useState<string>("");
  const [rowSelection, setRowSelection] = useState<Record<string, boolean>>({});
  const [isDeleting, setIsDeleting] = useState(false);

  const rows = useMemo(() => {
    return documents?.map(responseToDocument) || [];
  }, [documents]);

  const filteredData = useMemo(() => {
    if (!rows) return [];
    const q = query.trim().toLowerCase();
    return rows.filter((row) => {
      const matchesQuery = q ? row.title.toLowerCase().includes(q) : true;
      const matchesStatus = status ? row.documentStatus === status : true;
      const matchesType = docType ? row.documentType === docType : true;
      return matchesQuery && matchesStatus && matchesType;
    });
  }, [rows, query, status, docType]);

  const typeOptions = useMemo(() => {
    if (!rows) return [];
    return Array.from(new Set(rows.map((d) => d.documentType)));
  }, [rows]);

  const statusOptions = useMemo(() => {
    if (!rows) return [];
    return Array.from(new Set(rows.map((d) => d.documentStatus)));
  }, [rows]);

  useEffect(() => {
    const hasInProgress = rows.some(
      (d) => d.documentStatus !== "failed" && d.documentStatus !== "done",
    );
    if (!hasInProgress) return;
    const id = setInterval(() => {
      void mutate();
    }, 2000);
    return () => clearInterval(id);
  }, [rows, mutate]);

  const table = useReactTable({
    data: filteredData,
    columns: columnsWithSelection,
    getCoreRowModel: getCoreRowModel(),
    getRowId: (row) => row.id,
    onRowSelectionChange: setRowSelection,
    state: { rowSelection },
  });

  return (
    <div className="overflow-hidden rounded-md border">
      <div className="flex flex-wrap gap-2 p-3 border-b bg-muted/30">
        <div className="flex-1 min-w-[200px]">
          <Input
            placeholder="タイトルで検索"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
          />
        </div>
        <div className="flex items-center gap-2 overflow-x-auto whitespace-nowrap max-w-full">
          <select
            className="border-input bg-background text-foreground h-9 rounded-md border px-2 text-sm"
            value={docType}
            onChange={(e) => setDocType(e.target.value)}
          >
            <option value="">すべてのタイプ</option>
            {typeOptions.map((t) => (
              <option key={t} value={t}>
                {t}
              </option>
            ))}
          </select>
          <select
            className="border-input bg-background text-foreground h-9 rounded-md border px-2 text-sm"
            value={status}
            onChange={(e) => setStatus(e.target.value)}
          >
            <option value="">すべてのステータス</option>
            {statusOptions.map((s) => (
              <option key={s} value={s}>
                {s}
              </option>
            ))}
          </select>
          <Button
            variant="outline"
            onClick={async () => {
              const selected = table
                .getSelectedRowModel()
                .flatRows.map((r) => r.original);
              if (!selected.length) return;
              const ids = selected.map((s) => s.id);
              setIsDeleting(true);
              const results = await Promise.allSettled(
                ids.map((id) => deleteDocument(id)),
              );
              const failed = results.filter((r) => r.status === "rejected");
              if (failed.length === 0) {
                toast.success("削除しました");
              } else {
                toast.error(`${failed.length}件の削除に失敗しました`);
              }
              await mutate();
              setRowSelection({});
              setIsDeleting(false);
            }}
            disabled={
              isDeleting || table.getSelectedRowModel().flatRows.length === 0
            }
          >
            {isDeleting ? "削除中..." : "選択を削除"}
          </Button>
        </div>
      </div>
      <Table>
        <TableHeader>
          {table.getHeaderGroups().map((headerGroup) => (
            <TableRow key={headerGroup.id}>
              {headerGroup.headers.map((header) => {
                return (
                  <TableHead key={header.id}>
                    {header.isPlaceholder
                      ? null
                      : flexRender(
                          header.column.columnDef.header,
                          header.getContext(),
                        )}
                  </TableHead>
                );
              })}
            </TableRow>
          ))}
        </TableHeader>
        <TableBody>
          {table.getRowModel().rows?.length ? (
            table.getRowModel().rows.map((row) => (
              <TableRow
                key={row.id}
                data-state={row.getIsSelected() && "selected"}
              >
                {row.getVisibleCells().map((cell) => (
                  <TableCell key={cell.id}>
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </TableCell>
                ))}
              </TableRow>
            ))
          ) : (
            <TableRow>
              <TableCell
                colSpan={columnsWithSelection.length}
                className="h-full text-center"
              >
                <div className="flex items-center justify-center gap-2 text-muted-foreground">
                  {isLoading ? (
                    <>
                      <div className="animate-spin rounded-full h-4 w-4 border-2 border-primary border-t-transparent"></div>
                      読み込み中...
                    </>
                  ) : (
                    "結果がありません"
                  )}
                </div>
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </div>
  );
}
