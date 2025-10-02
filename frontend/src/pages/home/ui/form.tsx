"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import type { z } from "zod";
import { problemFormSchema } from "@/pages/home/model/zod";
import {
  Button,
  Form,
  FormControl,
  FormField,
  FormItem,
  FormMessage,
  Textarea,
} from "@/shared/ui";

interface ProblemFormProps {
  onSubmit: (values: z.infer<typeof problemFormSchema>) => void;
  isMutating: boolean;
}

export function ProblemForm({ onSubmit, isMutating }: ProblemFormProps) {
  const credentialsSignInForm = useForm<z.infer<typeof problemFormSchema>>({
    resolver: zodResolver(problemFormSchema),
    defaultValues: {
      description: "",
    },
  });

  return (
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
                  placeholder="解決したい課題について詳しく説明してください"
                  {...field}
                  className="resize-none flex-1 min-h-[60vh]"
                  disabled={isMutating}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <Button type="submit" className="w-full" disabled={isMutating}>
          {isMutating ? "作成中..." : "課題を作成する"}
        </Button>
      </form>
    </Form>
  );
}
