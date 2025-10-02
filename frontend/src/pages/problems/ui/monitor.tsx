import { LucideActivity } from "lucide-react";
import {
  Button,
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/shared/ui";
import type { Event } from "../model/zod";
import { EventList } from "./event-list";

type MonitorProps = {
  events: Event[];
};

export function Monitor({ events }: MonitorProps) {
  return (
    <div className="bg-white border rounded-lg p-4 space-y-3">
      {/* ヘッダー */}
      <div className="flex items-center space-x-2">
        <LucideActivity className="h-4 w-4 text-blue-500" />
        <span className="text-sm font-medium">処理中...</span>
      </div>

      {(() => {
        const actionEvents = events.filter(
          (event) => event.eventType === "action",
        );
        const lastActionEvent = actionEvents[actionEvents.length - 1];
        return (
          lastActionEvent && (
            <div className="text-sm text-gray-700">
              {lastActionEvent.message}
            </div>
          )
        );
      })()}

      <Sheet>
        <SheetTrigger className="w-full" asChild>
          <Button variant="outline" size="sm" className="w-full">
            詳細ログ
          </Button>
        </SheetTrigger>
        <SheetContent>
          <SheetHeader>
            <SheetTitle>詳細ログ</SheetTitle>
            <SheetDescription>詳細ログが表示されます。</SheetDescription>
          </SheetHeader>
          <EventList
            events={events.filter(
              (event, index, self) =>
                index === self.findIndex((e) => e.id === event.id),
            )}
          />
        </SheetContent>
      </Sheet>
    </div>
  );
}
