"use client";

import { useRouter } from "next/navigation";
import { toast } from "sonner";
import type { z } from "zod";
import type { problemFormSchema } from "@/pages/home/model/zod";
import { createHearing, executeHearing, useCreateProblem } from "@/shared/api";
import { Heading, RegularText } from "@/shared/ui";
import { ProblemForm } from "../ui/form";

export function HomePage() {
  const router = useRouter();
  const {
    trigger: createProblem,
    isMutating: isCreateProblemMutating,
    error: createProblemError,
  } = useCreateProblem();

  async function onSubmit(values: z.infer<typeof problemFormSchema>) {
    try {
      const problemResult = await createProblem({
        description: values.description,
      });

      if (createProblemError) {
        toast.error("Problem creation failed");
        return;
      }
      toast.success("Problem created successfully");
      router.push(`/problems/${problemResult.id}`);
      try {
        const hearingResult = await createHearing(problemResult.id);
        try {
          await executeHearing(problemResult.id, hearingResult.hearingId, {});
        } catch (_error) {
          toast.error("Hearing execution failed");
        }
      } catch (_error) {
        toast.error("Hearing creation failed");
        return;
      }
    } catch (_error) {
      toast.error("Problem creation failed");
    }
  }

  return (
    <div className="flex flex-col gap-5 h-full">
      <div className="flex gap-2 justify-between items-center">
        <Heading>Home</Heading>
      </div>
      <div className="flex flex-col items-center justify-center h-full">
        <Heading>Create Problem</Heading>
        <div className="flex w-full max-w-6xl h-full flex-col gap-5 p-5">
          <ProblemForm
            onSubmit={onSubmit}
            isMutating={isCreateProblemMutating}
          />
          <div className="flex items-center justify-center">
            <RegularText>課題を作成します。</RegularText>
          </div>
        </div>
      </div>
    </div>
  );
}
