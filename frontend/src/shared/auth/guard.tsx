"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { Loading } from "@/shared/ui";
import { useAuthContext } from "./use-auth-context";

export const AuthGuard = ({ children }: { children: React.ReactNode }) => {
  const { user, loading } = useAuthContext();
  const router = useRouter();
  useEffect(() => {
    if (!loading && !user) {
      router.push("/login");
    }
  }, [loading, user, router]);

  if (loading) {
    return <Loading />;
  }

  if (!user) {
    return <Loading />;
  }

  return children;
};
