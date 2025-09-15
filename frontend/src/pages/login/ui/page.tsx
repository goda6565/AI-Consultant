"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { signInWithEmailAndPassword } from "firebase/auth";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import type { z } from "zod";
import { credentialsSignInSchema } from "@/pages/login/model/zod";
import { auth } from "@/shared/config";
import {
  Button,
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
  Heading,
  Input,
  RegularText,
} from "@/shared/ui";

export function LoginPage() {
  const router = useRouter();

  const credentialsSignInForm = useForm<
    z.infer<typeof credentialsSignInSchema>
  >({
    resolver: zodResolver(credentialsSignInSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  });

  async function onSubmit(values: z.infer<typeof credentialsSignInSchema>) {
    await signInWithEmailAndPassword(auth, values.email, values.password)
      .then(() => {
        toast.success("Sign in successful");
        router.push("/");
      })
      .catch((error) => {
        toast.error(error.message);
      });
  }

  return (
    <div className="flex flex-col items-center justify-center h-screen">
      <Heading>Login</Heading>
      <div className="flex w-full max-w-lg flex-col gap-5 p-5">
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
        <div className="flex items-center justify-center">
          <RegularText>
            アカウントをお持ちでない方は管理者にお問い合わせください。
          </RegularText>
        </div>
      </div>
    </div>
  );
}
