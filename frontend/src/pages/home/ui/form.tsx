"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import type { z } from "zod";
import { problemFormSchema } from "@/pages/home/model/zod";
import { useCreateProblem } from "@/shared/api";
import {
  Button,
  Form,
  FormControl,
  FormField,
  FormItem,
  FormMessage,
  Heading,
  RegularText,
  Textarea,
} from "@/shared/ui";

export function ProblemForm() {
  const router = useRouter();
  const { trigger, isMutating, error } = useCreateProblem();

  const credentialsSignInForm = useForm<z.infer<typeof problemFormSchema>>({
    resolver: zodResolver(problemFormSchema),
    defaultValues: {
      description: "",
    },
  });

  async function onSubmit(values: z.infer<typeof problemFormSchema>) {
    const result = await trigger({ description: values.description });
    if (error) {
      toast.error(error.message as string);
    } else {
      toast.success("Problem created successfully");
      router.push(`/problems/${result.id}`);
    }
  }

  return (
    <div className="flex flex-col items-center justify-center h-full">
      <Heading>Create Problem</Heading>
      <div className="flex w-full max-w-6xl h-full flex-col gap-5 p-5">
        <Form {...credentialsSignInForm}>
          <form
            onSubmit={credentialsSignInForm.handleSubmit(onSubmit)}
            className="space-y-6 h-full flex flex-col"
          >
            <FormField
              control={credentialsSignInForm.control}
              name="description"
              render={({ field }) => (
                <FormItem className="flex-1 flex flex-col">
                  <FormControl>
                    <Textarea
                      placeholder="課題の詳細を入力してください"
                      {...field}
                      className="resize-none flex-1"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <Button type="submit" className="w-full" disabled={isMutating}>
              {isMutating ? "Creating..." : "Create Problem"}
            </Button>
          </form>
        </Form>
        <div className="flex items-center justify-center">
          <RegularText>課題を作成します。</RegularText>
        </div>
      </div>
    </div>
  );
}
