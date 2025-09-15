"use client";

import { useIsMobile } from "@/shared/hooks";
import { SidebarTrigger } from "@/shared/ui";

export function AppHeader() {
  const isMobile = useIsMobile();
  return (
    <>
      {isMobile && (
        <header className="h-10 w-full px-4 border-b flex items-center">
          <div className="flex items-center gap-2">
            <SidebarTrigger />
          </div>
        </header>
      )}
    </>
  );
}
