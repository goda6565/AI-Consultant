import { z } from "zod";

const envSchema = z.object({
  NEXT_PUBLIC_ADMIN_API_URL: z
    .string()
    .min(1, "NEXT_PUBLIC_ADMIN_API_URL is required"),
  NEXT_PUBLIC_AGENT_API_URL: z
    .string()
    .min(1, "NEXT_PUBLIC_AGENT_API_URL is required"),
});

export const env = envSchema.parse({
  NEXT_PUBLIC_ADMIN_API_URL: process.env.NEXT_PUBLIC_ADMIN_API_URL,
  NEXT_PUBLIC_AGENT_API_URL: process.env.NEXT_PUBLIC_AGENT_API_URL,
});
