import {
  LucideBrain,
  LucideCheck,
  LucideChevronDown,
  LucideDatabase,
  LucideEye,
  LucidePenTool,
  LucideSearch,
  LucideTarget,
} from "lucide-react";
import type React from "react";
import { useEffect, useRef, useState } from "react";
import {
  Button,
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  Item,
  ItemActions,
  ItemContent,
  ItemDescription,
  ItemMedia,
  ItemTitle,
  Markdown,
} from "@/shared/ui";
import type { Event } from "../model/zod";

export function EventList({ events }: { events: Event[] }) {
  const containerRef = useRef<HTMLDivElement | null>(null);

  // biome-ignore lint/correctness/useExhaustiveDependencies: scroll only when events changes
  useEffect(() => {
    const el = containerRef.current;
    if (!el) return;
    el.scrollTop = el.scrollHeight;
  }, [events]);

  return (
    <div
      ref={containerRef}
      className="max-h-full overflow-y-auto space-y-4 p-2 mb-4"
    >
      {events.map((event, _index) => (
        <div key={event.id} className="relative">
          {event.eventType === "action" ? (
            <ActionItem event={event} />
          ) : event.eventType === "input" ? (
            <InputItem event={event} />
          ) : (
            <OutputItem event={event} />
          )}
        </div>
      ))}
    </div>
  );
}

const ActionItem = ({ event }: { event: Event }) => {
  return (
    <Item variant="outline" size="sm">
      <ItemMedia>{actionTypeIconMap[event.actionType]}</ItemMedia>
      <ItemContent>
        <ItemTitle>{actionTypeLabels[event.actionType]}</ItemTitle>
        <ItemDescription>{event.message}</ItemDescription>
      </ItemContent>
    </Item>
  );
};

const InputItem = ({ event }: { event: Event }) => {
  return (
    <Item variant="outline" size="sm">
      <ItemMedia>{actionTypeIconMap[event.actionType]}</ItemMedia>
      <ItemContent>
        <ItemTitle>{actionTypeLabels[event.actionType]}</ItemTitle>
        <ItemDescription>{event.message}</ItemDescription>
      </ItemContent>
    </Item>
  );
};

const OutputItem = ({ event }: { event: Event }) => {
  const [isOpen, setIsOpen] = useState(false);
  return (
    <Item variant="outline" size="sm">
      <ItemMedia>{actionTypeIconMap[event.actionType]}</ItemMedia>
      <ItemContent className="min-w-0">
        <ItemTitle>{actionTypeLabels[event.actionType]}</ItemTitle>
        <ItemDescription>{event.message}</ItemDescription>
        <Dialog open={isOpen} onOpenChange={setIsOpen}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>{actionTypeLabels[event.actionType]}</DialogTitle>
              <DialogDescription>
                {actionTypeLabels[event.actionType]}の出力結果
              </DialogDescription>
            </DialogHeader>
            <div className="max-w-5xl max-h-[60vh] min-w-0 overflow-y-auto overflow-x-clip break-words">
              <Markdown>{event.message}</Markdown>
            </div>
          </DialogContent>
        </Dialog>
      </ItemContent>
      <ItemActions>
        <Button variant="ghost" size="sm" onClick={() => setIsOpen(!isOpen)}>
          <LucideChevronDown className="h-4 w-4" />
        </Button>
      </ItemActions>
    </Item>
  );
};

const actionTypeIconMap: Record<string, React.ReactNode> = {
  plan: <LucideTarget className="h-4 w-4 text-blue-600" />,
  externalSearch: <LucideSearch className="h-4 w-4 text-blue-600" />,
  internalSearch: <LucideDatabase className="h-4 w-4 text-blue-600" />,
  analyze: <LucideBrain className="h-4 w-4 text-blue-600" />,
  write: <LucidePenTool className="h-4 w-4 text-blue-600" />,
  review: <LucideEye className="h-4 w-4 text-blue-600" />,
  done: <LucideCheck className="h-4 w-4 text-blue-600" />,
};

const actionTypeLabels: Record<string, string> = {
  plan: "計画立案",
  externalSearch: "外部情報検索",
  internalSearch: "内部情報検索",
  analyze: "分析",
  write: "執筆",
  review: "レビュー",
  done: "完了",
};
