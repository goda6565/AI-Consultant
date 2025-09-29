import {
  LucideCheck,
  LucideEye,
  LucideFolderTree,
  LucideKeyboard,
  LucideMessageSquare,
  LucidePenTool,
  LucideSearch,
  LucideTarget,
  LucideZap,
} from "lucide-react";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
  ScrollArea,
} from "@/shared/ui";
import type { Event } from "../model/zod";

export function EventList({ events }: { events: Event[] }) {
  return (
    <ScrollArea className="w-full h-full">
      <div className="space-y-2 p-4 mb-16">
        {events.map((event) => (
          <div key={event.id}>
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
    </ScrollArea>
  );
}

const actionTypeIconMap: Record<string, React.ReactNode> = {
  plan: <LucideTarget className="h-4 w-4 text-purple-500" />,
  search: <LucideSearch className="h-4 w-4 text-blue-500" />,
  struct: <LucideFolderTree className="h-4 w-4 text-green-500" />,
  write: <LucidePenTool className="h-4 w-4 text-orange-500" />,
  review: <LucideEye className="h-4 w-4 text-teal-500" />,
  done: <LucideCheck className="h-4 w-4 text-emerald-500" />,
};

function ActionItem({ event }: { event: Event }) {
  return (
    <div className="flex items-start space-x-3 p-3 rounded-lg bg-card border shadow-sm hover:shadow-md transition-shadow duration-200">
      <div className="flex-shrink-0 mt-0.5">
        <div className="flex h-8 w-8 items-center justify-center rounded-full bg-muted relative z-10">
          {actionTypeIconMap[event.actionType] ?? (
            <LucideZap className="h-4 w-4 text-purple-500" />
          )}
        </div>
      </div>
      <div className="flex-1 min-w-0">
        <div className="flex items-center space-x-2 mb-1">
          <span className="text-xs font-medium text-muted-foreground uppercase tracking-wide">
            {event.actionType}
          </span>
          <div className="h-1 w-1 rounded-full bg-muted-foreground/50" />
          <span className="text-xs text-muted-foreground">Start</span>
        </div>
        <p className="text-sm text-black leading-relaxed line-clamp-2">
          {event.message}
        </p>
      </div>
    </div>
  );
}

function InputItem({ event }: { event: Event }) {
  return (
    <div className="flex items-start space-x-3 p-3 rounded-lg bg-card border shadow-sm hover:shadow-md transition-shadow duration-200">
      <div className="flex-shrink-0 mt-0.5">
        <div className="flex h-8 w-8 items-center justify-center rounded-full bg-muted relative z-10">
          <LucideKeyboard className="h-4 w-4 text-blue-500" />
        </div>
      </div>
      <div className="flex-1 min-w-0">
        <div className="flex items-center space-x-2 mb-1">
          <span className="text-xs font-medium text-muted-foreground uppercase tracking-wide">
            Input
          </span>
          <div className="h-1 w-1 rounded-full bg-muted-foreground/50" />
          <span className="text-xs text-muted-foreground">Input</span>
        </div>
        <p className="text-sm text-black leading-relaxed line-clamp-3">
          {event.message}
        </p>
      </div>
    </div>
  );
}

function OutputItem({ event }: { event: Event }) {
  return (
    <div className="rounded-lg bg-card border shadow-sm hover:shadow-md transition-shadow duration-200">
      <Accordion type="single" collapsible className="w-full">
        <AccordionItem value="item-1" className="border-none">
          <AccordionTrigger className="flex items-start space-x-3 p-3 text-left hover:no-underline">
            <div className="flex-shrink-0 mt-0.5">
              <div className="flex h-8 w-8 items-center justify-center rounded-full bg-muted relative z-10">
                <LucideMessageSquare className="h-4 w-4 text-green-500" />
              </div>
            </div>
            <div className="flex-1 min-w-0">
              <div className="flex items-center space-x-2 mb-1">
                <span className="text-xs font-medium text-muted-foreground uppercase tracking-wide">
                  Output
                </span>
                <div className="h-1 w-1 rounded-full bg-muted-foreground/50" />
                <span className="text-xs text-muted-foreground">Output</span>
              </div>
              <p className="text-sm text-black leading-relaxed line-clamp-1">
                {event.message}
              </p>
            </div>
          </AccordionTrigger>
          <AccordionContent className="px-3 pb-3">
            <div className="pl-11">
              <p className="text-sm text-black leading-relaxed whitespace-pre-wrap">
                {event.message}
              </p>
            </div>
          </AccordionContent>
        </AccordionItem>
      </Accordion>
    </div>
  );
}
