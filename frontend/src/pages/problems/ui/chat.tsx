"use client";

import { ArrowUpIcon } from "lucide-react";
import { useEffect, useRef, useState } from "react";
import ReactMarkdown from "react-markdown";
import { toast } from "sonner";
import type { Problem } from "@/shared/api";
import {
  useExecuteHearing,
  useGetHearing,
  useListHearingMessages,
} from "@/shared/api";
import { Button, Loading, ScrollArea, Textarea } from "@/shared/ui";
import type { Message } from "../model/zod";
import { Monitor } from "./monitor";

type ChatProps = {
  problem: Problem;
  mutateProblem: () => void;
};

export function Chat({ problem, mutateProblem }: ChatProps) {
  const [initialized, setInitialized] = useState(false);
  const [input, setInput] = useState("");
  const [localMessages, setLocalMessages] = useState<Message[]>([]);
  const scrollAreaRef = useRef<HTMLDivElement>(null);

  const { trigger: executeHearing, isMutating: isExecuteHearingMutating } =
    useExecuteHearing(problem.id);

  const { data: hearing, mutate: mutateHearing } = useGetHearing(problem.id, {
    swr: { enabled: problem.status !== "pending" },
  });

  const {
    data: hearingMessages,
    isLoading: isHearingMessagesLoading,
    mutate: mutateHearingMessages,
  } = useListHearingMessages(hearing?.id ?? "", {
    swr: { enabled: problem.status !== "pending" },
  });

  useEffect(() => {
    if (hearingMessages) {
      setLocalMessages(
        hearingMessages.hearingMessages.map((msg) => ({
          role: msg.role,
          message: msg.message,
        })),
      );
    }
  }, [hearingMessages]);

  useEffect(() => {
    if (scrollAreaRef.current) {
      const scrollContainer = scrollAreaRef.current.querySelector(
        "[data-radix-scroll-area-viewport]",
      );
      if (scrollContainer) {
        scrollContainer.scrollTop = scrollContainer.scrollHeight;
      }
    }
  });

  // hearing 初期化
  useEffect(() => {
    if (problem.status === "pending" && !initialized) {
      setInitialized(true);
      executeHearing({ user_message: null })
        .then(() => {
          mutateHearing();
          mutateProblem();
          mutateHearingMessages();
        })
        .catch((err) =>
          toast.error(err.message || "Hearing の初期化に失敗しました"),
        );
    }
  }, [
    problem.status,
    executeHearing,
    mutateHearing,
    mutateHearingMessages,
    mutateProblem,
    initialized,
  ]);

  if (problem.status === "pending")
    return (
      <div className="h-full flex items-center justify-center">
        <Loading />
      </div>
    );

  return (
    <div className="flex flex-col bg-white h-full max-w-4xl mx-auto">
      {/* list of messages */}
      <div className="flex-1 min-h-0">
        <ScrollArea ref={scrollAreaRef} className="h-full">
          <div className="p-4 space-y-6">
            {localMessages.map((msg, index) => (
              <div
                key={`${msg.role}-${index}`}
                className={`flex ${
                  msg.role === "user" ? "justify-end" : "justify-start"
                }`}
              >
                <div
                  className={`${
                    msg.role === "user"
                      ? "px-4 py-3 rounded-2xl shadow-sm bg-gray-100 text-black rounded-br-md max-w-[70%]"
                      : "w-full"
                  }`}
                >
                  {msg.role === "user" ? (
                    <div className="text-sm leading-relaxed whitespace-pre-wrap">
                      {msg.message}
                    </div>
                  ) : (
                    <div className="leading-relaxed prose prose-sm w-full text-black max-w-none">
                      <ReactMarkdown>{msg.message}</ReactMarkdown>
                    </div>
                  )}
                </div>
              </div>
            ))}
          </div>
          {problem.status === "processing" && <Monitor problem={problem} />}
        </ScrollArea>
      </div>

      {/* 入力欄 */}
      <form
        onSubmit={async (e) => {
          e.preventDefault();
          if (!input.trim()) return;

          const newMessage: Message = {
            role: "user",
            message: input,
          };
          setLocalMessages((prev) => [...prev, newMessage as Message]);
          setInput("");

          try {
            await executeHearing({ user_message: input });
            mutateProblem();
            mutateHearingMessages(); // バックグラウンド同期
          } catch (_err) {
            toast.error("メッセージ送信に失敗しました");
          }
        }}
        className="sticky bottom-0 bg-white flex items-center gap-2 p-4"
      >
        <Textarea
          placeholder="メッセージを入力... (Enterで送信、Shift+Enterで改行)"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === "Enter" && !e.shiftKey) {
              e.preventDefault();
              if (
                input.trim() &&
                !isExecuteHearingMutating &&
                !isHearingMessagesLoading &&
                problem.status !== "processing"
              ) {
                e.currentTarget.form?.requestSubmit();
              }
            }
          }}
          disabled={
            isHearingMessagesLoading ||
            isExecuteHearingMutating ||
            problem.status === "processing"
          }
          className="flex-1 min-h-[40px] max-h-[240px] resize-none rounded-4xl p-4"
          rows={1}
        />
        <Button
          type="submit"
          disabled={
            !input.trim() ||
            isExecuteHearingMutating ||
            isHearingMessagesLoading ||
            problem.status === "processing"
          }
          className="rounded-full w-[40px] h-[40px] p-0 flex items-center justify-center"
        >
          {isExecuteHearingMutating ? (
            <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
          ) : (
            <ArrowUpIcon className="w-4 h-4" />
          )}
        </Button>
      </form>
    </div>
  );
}
