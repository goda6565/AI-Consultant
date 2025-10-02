import { zodResolver } from "@hookform/resolvers/zod";
import { LucideArrowUp, LucideLoader2 } from "lucide-react";
import { useForm } from "react-hook-form";
import type { z } from "zod";
import {
  Button,
  Form,
  FormControl,
  FormField,
  FormItem,
  FormMessage,
  Textarea,
} from "@/shared/ui";
import { MessageFormSchema } from "../model/zod";

interface MessageFormProps {
  onSubmit: (values: z.infer<typeof MessageFormSchema>) => void | Promise<void>;
  isMutating: boolean;
}

export function MessageForm({ onSubmit, isMutating }: MessageFormProps) {
  const messageForm = useForm<z.infer<typeof MessageFormSchema>>({
    resolver: zodResolver(MessageFormSchema),
    defaultValues: {
      message: "",
    },
  });
  const handleSubmit = messageForm.handleSubmit(async (values) => {
    messageForm.reset({ message: "" });
    await onSubmit(values);
  });
  return (
    <Form {...messageForm}>
      <form onSubmit={handleSubmit} className="space-x-6 flex items-center">
        <FormField
          control={messageForm.control}
          name="message"
          render={({ field }) => (
            <FormItem className="flex-1 flex flex-col max-h-64">
              <FormControl>
                <Textarea
                  placeholder="メッセージを入力してください"
                  {...field}
                  className="resize-none flex-1"
                  disabled={isMutating}
                  onKeyDown={(e) => {
                    if (e.key === "Enter" && !e.shiftKey) {
                      e.preventDefault();
                      handleSubmit();
                    }
                  }}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <Button
          type="submit"
          className="w-10 h-10 rounded-full"
          disabled={isMutating}
        >
          {isMutating ? (
            <LucideLoader2 className="w-full h-full animate-spin" />
          ) : (
            <LucideArrowUp className="w-full h-full" />
          )}
        </Button>
      </form>
    </Form>
  );
}
