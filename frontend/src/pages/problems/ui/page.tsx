"use client";

import { redirect } from "next/navigation";
import { use } from "react";
import { toast } from "sonner";
import type { z } from "zod";
import { useExecuteHearing } from "@/shared/api";
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
import type { Message, MessageFormSchema } from "../model/zod";
import { MessageForm } from "./form";
import { MessageView } from "./message-view";
import { Monitor } from "./monitor";

type ProblemPageProps = {
  params: Promise<{ id: string }>;
};

export function ProblemPage({ params }: ProblemPageProps) {
  const { id } = use(params);
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

  if (isChatError || isHearingMessagesError) {
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
      <div className="flex flex-col flex-1 justify-between min-h-0 max-w-3xl w-full mx-auto">
        <div className="flex-1 min-h-0">
          <ScrollArea className="h-full">
            <div className="mb-30">
              <MessageView messages={localMessages} />
              {problem.status === "processing" && <Monitor events={events} />}
            </div>
          </ScrollArea>
        </div>
        <MessageForm
          onSubmit={onChatSubmit}
          isMutating={isExecuteHearingMutating}
        />
      </div>
    </div>
  );
}
