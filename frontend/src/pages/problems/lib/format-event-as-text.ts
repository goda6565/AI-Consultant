import type { Event } from "../model/zod";

export const formatEventsAsText = (events: Event[]): string => {
  return events
    .map((event, index) => {
      const eventTypeLabel =
        event.eventType === "action"
          ? "アクション"
          : event.eventType === "input"
            ? "入力"
            : "出力";
      const actionTypeLabel = event.actionType ? ` (${event.actionType})` : "";

      return `${index + 1}. [${eventTypeLabel}${actionTypeLabel}] ${event.message}`;
    })
    .join("\n\n");
};
