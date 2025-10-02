"use client";

import { signInWithEmailAndPassword } from "firebase/auth";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import type { z } from "zod";
import type { credentialsSignInSchema } from "@/pages/login/model/zod";
import { LoginForm } from "@/pages/login/ui/form";
import { auth } from "@/shared/config";
import { Heading, RegularText } from "@/shared/ui";

export function LoginPage() {
  const router = useRouter();

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
        <LoginForm onSubmit={onSubmit} />
        <div className="flex items-center justify-center">
          <RegularText>
            アカウントをお持ちでない方は管理者にお問い合わせください。
          </RegularText>
        </div>
      </div>
    </div>
  );
}
