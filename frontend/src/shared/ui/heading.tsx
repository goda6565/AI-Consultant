"use client";

import { useIsMobile } from "@/shared/hooks";

export function Heading({ children }: { children: React.ReactNode }) {
  const isMobile = useIsMobile();
  return (
    <h1 className={`${isMobile ? "text-xl" : "text-2xl"} font-bold`}>
      {children}
    </h1>
  );
}
