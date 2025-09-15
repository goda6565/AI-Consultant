import { z } from "zod";

const envSchema = z.object({
  NEXT_PUBLIC_ADMIN_API_URL: z
    .string()
    .min(1, "NEXT_PUBLIC_ADMIN_API_URL is required"),
});

// Note: In Next.js, only direct references to process.env.NEXT_PUBLIC_* are inlined on the client.
// Do NOT pass process.env wholesale; instead reference keys explicitly.
export const env = envSchema.parse({
  NEXT_PUBLIC_ADMIN_API_URL: process.env.NEXT_PUBLIC_ADMIN_API_URL,
});
