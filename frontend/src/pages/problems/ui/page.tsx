"use client";

import { redirect } from "next/navigation";
import { use, useEffect, useRef } from "react";
import { toast } from "sonner";
import type { z } from "zod";
import { useExecuteHearing, useGetReport } from "@/shared/api";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
  Badge,
  Heading,
  Loading,
  ScrollArea,
} from "@/shared/ui";
import { useChatApi } from "../api/use-chat-api";
import { useChatMessage } from "../api/use-chat-message";
import { useEventSse } from "../api/use-event";
import { formatEventsAsText } from "../lib/format-event-as-text";
import type { Message, MessageFormSchema } from "../model/zod";
import { MessageForm } from "./form";
import { MessageView } from "./message-view";
import { Monitor } from "./monitor";
import { ReportView } from "./report-view";

type ProblemPageProps = {
  params: Promise<{ id: string }>;
};

export function ProblemPage({ params }: ProblemPageProps) {
  const { id } = use(params);
  const scrollRef = useRef<HTMLDivElement>(null);
  // Fetch API
  const { problem, hearing, mutateChat, isChatLoading, isChatError } =
    useChatApi(id);

  const { events } = useEventSse({
    problemId: id,
    enabled: problem?.status === "processing",
  });

  const {
    localMessages,
    setLocalMessages,
    isLoading: isHearingMessagesLoading,
    error: isHearingMessagesError,
  } = useChatMessage({
    hearingId: hearing?.id ?? "",
    enabled: problem?.status !== "pending",
  });

  const {
    data: report,
    isLoading: isReportLoading,
    error: isReportError,
  } = useGetReport(id, {
    swr: {
      enabled: problem?.status === "done",
    },
  });

  // biome-ignore lint/correctness/useExhaustiveDependencies: scroll only when length changes
  useEffect(() => {
    if (scrollRef.current) {
      // ScrollAreaの内部コンテンツを最下部にスクロール
      const scrollContainer = scrollRef.current.querySelector(
        "[data-radix-scroll-area-viewport]",
      );
      if (scrollContainer) {
        scrollContainer.scrollTop = scrollContainer.scrollHeight;
      }
    }
  }, [localMessages.length, !!report, events.length]);

  // Mutation API
  const { trigger: executeHearing, isMutating: isExecuteHearingMutating } =
    useExecuteHearing(problem?.id ?? "", hearing?.id ?? "");

  if (isChatLoading || isHearingMessagesLoading) {
    return (
      <div className="flex items-center justify-center h-full">
        <Loading />
      </div>
    );
  }

  if (isChatError || isHearingMessagesError || isReportError) {
    toast.error("failed to fetch messages");
  }

  if (!problem) {
    redirect("/");
  }

  if (localMessages.length === 0) {
    return (
      <div className="flex items-center justify-center h-full">
        <Loading />
      </div>
    );
  }

  const onChatSubmit = async (values: z.infer<typeof MessageFormSchema>) => {
    const newMessage: Message = {
      role: "user",
      message: values.message,
    };
    setLocalMessages((prev) => [...prev, newMessage]);
    try {
      const result = await executeHearing({ user_message: values.message });
      const newMessage: Message = {
        role: "assistant",
        message: result.assistant_message,
      };
      setLocalMessages((prev) => [...prev, newMessage]);
      mutateChat();
    } catch (_err) {
      toast.error("メッセージ送信に失敗しました");
    }
  };

  const onCopyEvents = () => {
    try {
      const text = formatEventsAsText(events);
      navigator.clipboard.writeText(text);
      toast.success("イベントログをコピーしました");
    } catch (_err) {
      toast.error("イベントログをコピーに失敗しました");
    }
  };

  const onCopyReport = () => {
    try {
      navigator.clipboard.writeText(report?.content ?? "");
      toast.success("レポートをコピーしました");
    } catch (_err) {
      toast.error("レポートをコピーに失敗しました");
    }
  };

  return (
    <div className="flex flex-col h-full">
      <div className="flex gap-6 justify-between items-center p-4 border-b">
        <Accordion type="single" collapsible className="w-full">
          <AccordionItem value="item-1">
            <AccordionTrigger>
              <Heading>{problem.title}</Heading>
            </AccordionTrigger>
            <AccordionContent>{problem.description}</AccordionContent>
          </AccordionItem>
        </Accordion>
        <Badge variant="outline">{problem.status}</Badge>
      </div>
      <div className="flex flex-col flex-1 justify-between min-h-0 max-w-4xl w-full mx-auto">
        <div className="flex-1 min-h-0">
          <ScrollArea className="h-full" ref={scrollRef}>
            <div className="mb-30">
              <MessageView messages={localMessages} />
              {problem.status === "processing" && (
                <Monitor events={events} onCopyEvents={onCopyEvents} />
              )}
              {problem.status === "done" && (
                <ReportView
                  report={report}
                  isLoading={isReportLoading}
                  onCopyReport={onCopyReport}
                />
              )}
            </div>
          </ScrollArea>
        </div>
        <MessageForm
          onSubmit={onChatSubmit}
          isMutating={isExecuteHearingMutating}
          enabled={problem.status === "hearing"}
        />
      </div>
    </div>
  );
}
