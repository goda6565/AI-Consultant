import {
  LucideBrain,
  LucideCheck,
  LucideEye,
  LucideKeyboard,
  LucidePenTool,
  LucideSearch,
  LucideTarget,
  LucideZap,
} from "lucide-react";
import type React from "react";
import { useEffect, useRef } from "react";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
  Badge,
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
    <div className="border-l-4 border-l-blue-500 bg-blue-50/50 rounded-r-lg p-4">
      <div className="flex items-center gap-3">
        <div className="flex h-6 w-6 items-center justify-center rounded-full bg-blue-100 flex-shrink-0">
          {actionTypeIconMap[event.actionType]}
        </div>
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2 mb-1">
            <span className="text-sm font-medium">
              {actionTypeLabels[event.actionType]}
            </span>
            <Badge variant="secondary" className="text-xs px-1.5 py-0.5">
              {event.actionType}
            </Badge>
          </div>
          <p className="text-sm text-muted-foreground line-clamp-2 whitespace-pre-wrap">
            {event.message}
          </p>
        </div>
      </div>
    </div>
  );
};

const InputItem = ({ event }: { event: Event }) => {
  return (
    <Accordion type="single" collapsible className="w-full">
      <AccordionItem
        value="input"
        className="border-l-4 border-l-green-500 bg-green-50/50 rounded-r-lg"
      >
        <AccordionTrigger className="px-4 py-3 hover:no-underline">
          <div className="flex items-center gap-3">
            <div className="flex h-6 w-6 items-center justify-center rounded-full bg-green-100 flex-shrink-0">
              <LucideKeyboard className="h-4 w-4 text-green-600" />
            </div>
            <div className="flex-1 min-w-0 text-left">
              <div className="flex items-center gap-2">
                <span className="text-sm font-medium">入力</span>
                <Badge variant="outline" className="text-xs px-1.5 py-0.5">
                  {event.actionType}
                </Badge>
              </div>
            </div>
          </div>
        </AccordionTrigger>
        <AccordionContent className="px-4 pb-4">
          <p className="text-sm text-muted-foreground line-clamp-5 whitespace-pre-wrap">
            {event.message}
          </p>
        </AccordionContent>
      </AccordionItem>
    </Accordion>
  );
};

const OutputItem = ({ event }: { event: Event }) => {
  return (
    <Accordion type="single" collapsible className="w-full">
      <AccordionItem
        value="output"
        className="border-l-4 border-l-purple-500 bg-purple-50/50 rounded-r-lg"
      >
        <AccordionTrigger className="px-4 py-3 hover:no-underline">
          <div className="flex items-center gap-3">
            <div className="flex h-6 w-6 items-center justify-center rounded-full bg-purple-100 flex-shrink-0">
              <LucideZap className="h-4 w-4 text-purple-600" />
            </div>
            <div className="flex-1 min-w-0 text-left">
              <div className="flex items-center gap-2">
                <span className="text-sm font-medium">出力</span>
                <Badge variant="outline" className="text-xs px-1.5 py-0.5">
                  {event.actionType}
                </Badge>
              </div>
            </div>
          </div>
        </AccordionTrigger>
        <AccordionContent className="px-4 pb-4">
          <p className="text-sm text-muted-foreground line-clamp-10 whitespace-pre-wrap">
            {event.message}
          </p>
        </AccordionContent>
      </AccordionItem>
    </Accordion>
  );
};

const actionTypeIconMap: Record<string, React.ReactNode> = {
  plan: <LucideTarget className="h-4 w-4 text-blue-600" />,
  search: <LucideSearch className="h-4 w-4 text-blue-600" />,
  analyze: <LucideBrain className="h-4 w-4 text-blue-600" />,
  write: <LucidePenTool className="h-4 w-4 text-blue-600" />,
  review: <LucideEye className="h-4 w-4 text-blue-600" />,
  done: <LucideCheck className="h-4 w-4 text-blue-600" />,
};

const actionTypeLabels: Record<string, string> = {
  plan: "計画立案",
  search: "情報検索",
  analyze: "分析",
  write: "執筆",
  review: "レビュー",
  done: "完了",
};
