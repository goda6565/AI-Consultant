"use client";

import { useEffect, useState } from "react";
import { toast } from "sonner";
import type { Problem } from "@/shared/api";
import {
  useExecuteHearing,
  useGetHearing,
  useListHearingMessages,
} from "@/shared/api";
import { Button, Loading, ScrollArea, Textarea } from "@/shared/ui";

type ChatProps = {
  problem: Problem;
  mutateProblem: () => void;
};

export function Chat({ problem, mutateProblem }: ChatProps) {
  const [initialized, setInitialized] = useState(false);
  const [input, setInput] = useState("");

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

  // 初回読み込み時はローディング表示
  if (problem.status === "pending") return <Loading />;

  return (
    <div className="flex flex-col bg-white h-full">
      {/* list of messages */}
      <div className="flex-1 min-h-0">
        <ScrollArea className="h-full">
          <div className="p-4 space-y-4">
            {hearingMessages?.hearingMessages.map((msg) => (
              <div
                key={msg.id}
                className={`flex ${
                  msg.role === "user" ? "justify-end" : "justify-start"
                }`}
              >
                <div
                  className={`max-w-[75%] px-4 py-3 rounded-2xl shadow-sm ${
                    msg.role === "user"
                      ? "bg-gray-100 text-gray-900 rounded-br-md"
                      : "bg-white text-gray-900 rounded-bl-md border"
                  }`}
                >
                  <div className="text-sm leading-relaxed whitespace-pre-wrap">
                    {msg.message}
                  </div>
                </div>
              </div>
            ))}
          </div>
        </ScrollArea>
      </div>

      {/* 入力欄 */}
      <form
        onSubmit={async (e) => {
          e.preventDefault();
          if (!input.trim()) return;
          try {
            await executeHearing({ user_message: input });
            await mutateHearingMessages();
            mutateProblem();
            setInput("");
          } catch (_err) {
            toast.error("メッセージ送信に失敗しました");
          }
        }}
        className="sticky bottom-0 bg-white flex items-center gap-2 p-4 border-t"
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
          className="flex-1 min-h-[40px] max-h-[120px] resize-none"
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
        >
          {isExecuteHearingMutating ? (
            <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
          ) : (
            "送信"
          )}
        </Button>
      </form>
    </div>
  );
}
