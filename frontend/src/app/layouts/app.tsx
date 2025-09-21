import { cookies } from "next/headers";
import { AuthProvider } from "@/app/provider/auth";
import { SidebarProvider } from "@/shared/ui/sidebar";
import { AppHeader } from "@/widgets/app-header";
import { AppSidebar } from "@/widgets/app-sidebar";

export async function AppLayout({ children }: { children: React.ReactNode }) {
  const cookieStore = await cookies();
  const defaultOpen = cookieStore.get("sidebar_state")?.value === "true";
  return (
    <AuthProvider>
      <SidebarProvider defaultOpen={defaultOpen}>
        <AppSidebar />
        <main className="flex flex-col w-full h-screen">
          <AppHeader />
          <div className="flex-1 min-h-0 p-4">{children}</div>
        </main>
      </SidebarProvider>
    </AuthProvider>
  );
}
