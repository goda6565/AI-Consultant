import { toast } from "sonner";
import type { Report } from "@/shared/api";
import { Button, Loading, Markdown } from "@/shared/ui";

type ReportViewProps = {
  report?: Report;
  isLoading: boolean;
  onCopyReport: () => void;
};

export function ReportView({
  report,
  isLoading,
  onCopyReport,
}: ReportViewProps) {
  if (isLoading) {
    return <Loading />;
  }

  if (!report) {
    toast.error("Report not found");
    return null;
  }

  return (
    <div className="p-8 rounded-xl border">
      <div className="flex justify-end items-center">
        <Button variant="outline" onClick={onCopyReport}>
          Copy
        </Button>
      </div>
      <div className="mt-2 prose prose-sm max-w-none">
        <Markdown>{report.content}</Markdown>
      </div>
    </div>
  );
}
