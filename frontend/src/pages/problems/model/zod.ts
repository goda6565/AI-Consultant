import { z } from "zod";

export const MessageFormSchema = z.object({
  message: z.string(),
});

export const MessageSchema = z.object({
  role: z.enum(["user", "assistant"]),
  message: z.string(),
});

export type Message = z.infer<typeof MessageSchema>;

export const EventSchema = z.object({
  id: z.string(),
  eventType: z.enum(["action", "input", "output"]),
  actionType: z.enum(["plan", "search", "analyze", "write", "review", "done"]),
  message: z.string(),
});

export type Event = z.infer<typeof EventSchema>;
