import { z } from "zod";

export const problemFormSchema = z.object({
  description: z.string().min(1),
});
