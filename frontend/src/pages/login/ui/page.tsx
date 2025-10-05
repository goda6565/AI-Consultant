"use client";

import { signInWithEmailAndPassword } from "firebase/auth";
import Image from "next/image";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import type { z } from "zod";
import type { credentialsSignInSchema } from "@/pages/login/model/zod";
import { LoginForm } from "@/pages/login/ui/form";
import { auth } from "@/shared/config";

export function LoginPage() {
  const router = useRouter();

  async function onSubmit(values: z.infer<typeof credentialsSignInSchema>) {
    await signInWithEmailAndPassword(auth, values.email, values.password)
      .then(() => {
        toast.success("Sign in successful");
        router.push("/");
      })
      .catch(() => {
        toast.error("Password or email is incorrect");
      });
  }

  return (
    <div className="grid min-h-svh lg:grid-cols-2">
      <div className="flex flex-col gap-4 p-6 md:p-10">
        <div className="flex justify-center gap-2 md:justify-start">
          <div className="flex items-center justify-center font-bold">
            <div className="text-primary-foreground flex items-center justify-center mb-2">
              <Image src="/logo.png" alt="Logo" width={60} height={60} />
            </div>
            <span className="-ml-2">Consultant</span>
          </div>
        </div>
        <div className="flex flex-1 items-center justify-center">
          <div className="w-full max-w-md mb-14">
            <LoginForm onSubmit={onSubmit} />
          </div>
        </div>
      </div>
      <div className="bg-muted relative hidden lg:block">
        <Image
          width={1000}
          height={1000}
          src="/login.jpg"
          alt="Image"
          className="absolute inset-0 h-full w-full object-cover dark:brightness-[0.2] dark:grayscale"
        />
      </div>
    </div>
  );
}
