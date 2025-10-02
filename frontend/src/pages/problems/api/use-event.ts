import { useEffect, useState } from "react";
import { toast } from "sonner";
import { useListEvents } from "@/shared/api";
import { env } from "@/shared/config";
import type { Event } from "../model/zod";
import { EventSchema } from "../model/zod";

type UseEventSseProps = {
  problemId: string;
  enabled: boolean;
};

export const useEventSse = ({ problemId, enabled }: UseEventSseProps) => {
  const { data: eventsData, isLoading, error } = useListEvents(problemId);

  const [events, setEvents] = useState<Event[]>([]);

  useEffect(() => {
    if (eventsData) {
      setEvents(eventsData.events);
    }
  }, [eventsData]);

  useEffect(() => {
    if (!enabled) return;
    let es: EventSource | null = null;

    const connect = () => {
      es = new EventSource(
        `${env.NEXT_PUBLIC_ADMIN_API_URL}/api/events/${problemId}/stream`,
      );

      es.onmessage = (e) => {
        try {
          const parsed = EventSchema.parse(JSON.parse(e.data));
          setEvents((prev) => [...prev, parsed]);
        } catch {
          toast.error("failed to parse event");
        }
      };

      es.onerror = () => {
        toast.error("event stream connection closed. reconnecting...");
        es?.close();
        setTimeout(connect, 3000);
      };
    };

    connect();

    return () => {
      es?.close();
    };
  }, [problemId, enabled]);

  return { events, isLoading, error };
};
