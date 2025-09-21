import { z } from "zod";

export const requestMessageSchema = z.object({
  user_message: z.string(),
});

export const responseMessageSchema = z.object({
  type: z.enum(["hearing_response", "hearing_completed", "error"]),
  assistant_message: z.string(),
});

export const MessageSchema = z.object({
  id: z.string(),
  content: z.string(),
  isUser: z.boolean(),
  timestamp: z.date(),
});

export type Message = z.infer<typeof MessageSchema>;
