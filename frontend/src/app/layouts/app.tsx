import { AuthGuard } from "@/shared/auth";
import { SidebarProvider } from "@/shared/ui/sidebar";
import { AppHeader } from "@/widgets/app-header";
import { AppSidebar } from "@/widgets/app-sidebar";

export function AppLayout({ children }: { children: React.ReactNode }) {
  return (
    <AuthGuard>
      <SidebarProvider>
        <AppSidebar />
        <main className="flex flex-col w-full h-screen">
          <AppHeader />
          <div className="flex-1 min-h-0 p-4">{children}</div>
        </main>
      </SidebarProvider>
    </AuthGuard>
  );
}
