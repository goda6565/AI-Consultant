import { Markdown } from "@/shared/ui";
import type { Message } from "../model/zod";

type MessageViewProps = {
  messages: Message[];
};

export function MessageView({ messages }: MessageViewProps) {
  return (
    <div className="p-4 space-y-12">
      {messages.map((msg, index) => (
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
                <Markdown>{msg.message}</Markdown>
              </div>
            )}
          </div>
        </div>
      ))}
    </div>
  );
}
