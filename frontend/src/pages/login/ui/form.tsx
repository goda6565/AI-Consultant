import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import type { z } from "zod";
import { credentialsSignInSchema } from "@/pages/login/model/zod";
import {
  Button,
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
  Input,
} from "@/shared/ui";

interface LoginFormProps {
  onSubmit: (values: z.infer<typeof credentialsSignInSchema>) => void;
}

export function LoginForm({ onSubmit }: LoginFormProps) {
  const credentialsSignInForm = useForm<
    z.infer<typeof credentialsSignInSchema>
  >({
    resolver: zodResolver(credentialsSignInSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  });
  return (
    <Form {...credentialsSignInForm}>
      <form
        onSubmit={credentialsSignInForm.handleSubmit(onSubmit)}
        className="space-y-6"
      >
        <FormField
          control={credentialsSignInForm.control}
          name="email"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Email</FormLabel>
              <FormControl>
                <Input type="email" placeholder="Email" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={credentialsSignInForm.control}
          name="password"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Password</FormLabel>
              <FormControl>
                <Input type="password" placeholder="Password" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <Button type="submit" className="w-full">
          Sign In
        </Button>
      </form>
    </Form>
  );
}
