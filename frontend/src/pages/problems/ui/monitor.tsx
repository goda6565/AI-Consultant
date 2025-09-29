import { LucideActivity } from "lucide-react";
import { useEffect, useState } from "react";
import { toast } from "sonner";
import { type Problem, useListEvents } from "@/shared/api";
import { env } from "@/shared/config";
import {
  Button,
  Drawer,
  DrawerContent,
  DrawerTitle,
  DrawerTrigger,
  Loading,
} from "@/shared/ui";
import type { Event } from "../model/zod";
import { EventList } from "./event-list";

type MonitorProps = {
  problem: Problem;
};

export function Monitor({ problem }: MonitorProps) {
  const [events, setEvents] = useState<Event[]>([]);

  const {
    data: eventsData,
    isLoading,
    error,
  } = useListEvents("cc2546eb-1d32-4697-beda-b3ff2ab9bb92");

  useEffect(() => {
    if (eventsData) {
      setEvents(eventsData.events);
    }
  }, [eventsData]);

  useEffect(() => {
    const url = `${env.NEXT_PUBLIC_ADMIN_API_URL}/api/events/cc2546eb-1d32-4697-beda-b3ff2ab9bb92/stream`;
    const es = new EventSource(url);

    es.onmessage = (e) => {
      try {
        const parsed = JSON.parse(e.data) as Event;
        setEvents((prev) => [...prev, parsed]);
      } catch (_err) {
        toast.error("イベントストリームのパースに失敗しました");
      }
    };

    es.onerror = (_err) => {
      toast.error("イベントストリームの接続に失敗しました");
      es.close();
    };

    return () => {
      es.close();
    };
  }, []);

  if (isLoading)
    return (
      <div className="flex items-center justify-center h-full">
        <Loading />
      </div>
    );
  if (error) toast.error(error.message);

  return (
    <div className="bg-white border rounded-lg p-4 space-y-3">
      {/* ヘッダー */}
      <div className="flex items-center space-x-2">
        <LucideActivity className="h-4 w-4 text-blue-500" />
        <span className="text-sm font-medium">処理中...</span>
      </div>

      {/* 現在のステップ */}
      {events[events.length - 1] && (
        <div className="text-sm text-gray-700">
          {events[events.length - 1].message}
        </div>
      )}

      {/* 詳細表示ボタン */}
      <Drawer direction="right">
        <DrawerTrigger asChild>
          <Button variant="outline" size="sm" className="w-full">
            詳細ログ
          </Button>
        </DrawerTrigger>
        <DrawerContent>
          <DrawerTitle className="text-center p-4">{problem.title}</DrawerTitle>
          <EventList events={events} />
        </DrawerContent>
      </Drawer>
    </div>
  );
}
