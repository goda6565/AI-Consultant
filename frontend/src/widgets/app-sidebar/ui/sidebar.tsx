"use client";

import { signOut } from "firebase/auth";
import {
  ChevronUp,
  Ellipsis,
  FileText,
  Plus,
  Search,
  User2,
} from "lucide-react";
import Image from "next/image";
import { useEffect, useState } from "react";
import { toast } from "sonner";
import { deleteProblem, useListProblems } from "@/shared/api";
import { auth } from "@/shared/config";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
  Logo,
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuAction,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarTrigger,
  useSidebar,
} from "@/shared/ui";

const menuItems = [
  {
    title: "新しい課題",
    url: "/",
    icon: Plus,
  },
  {
    title: "課題を検索",
    url: "/cases",
    icon: Search,
  },
  {
    title: "ドキュメントを追加",
    url: "/documents",
    icon: FileText,
  },
];

export function AppSidebar() {
  const { open } = useSidebar();
  const [isHovering, setIsHovering] = useState(false);
  const user = auth.currentUser;
  const { data: problems, mutate: mutateProblems } = useListProblems();

  useEffect(() => {
    if (open) {
      setIsHovering(false);
    }
  }, [open]);

  const handleMouseEnter = () => {
    if (!open) {
      setIsHovering(true);
    }
  };

  const handleMouseLeave = () => {
    setIsHovering(false);
  };

  return (
    <Sidebar
      collapsible="icon"
      onMouseEnter={handleMouseEnter}
      onMouseLeave={handleMouseLeave}
    >
      <SidebarHeader>
        <div className="flex items-center gap-2 h-10">
          {open ? (
            <div className="flex items-center justify-between w-full">
              <Logo />
              <SidebarTrigger />
            </div>
          ) : (
            <div className="relative w-[30px] h-[30px]">
              <Logo
                className={`absolute inset-0 transition-opacity duration-200 ${isHovering ? "opacity-0" : "opacity-100"}`}
              />
              <div
                className={`absolute inset-0 transition-opacity duration-200 ${isHovering ? "opacity-100" : "opacity-0"}`}
              >
                <SidebarTrigger />
              </div>
            </div>
          )}
        </div>
      </SidebarHeader>
      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupContent>
            <SidebarMenu>
              {menuItems.map((item) => (
                <SidebarMenuItem key={item.title}>
                  <SidebarMenuButton asChild>
                    <a href={item.url}>
                      <item.icon />
                      <span>{item.title}</span>
                    </a>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
        {open && (
          <SidebarGroup>
            <SidebarGroupLabel>課題</SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                {problems?.problems.map((item) => (
                  <SidebarMenuItem key={item.id}>
                    <SidebarMenuButton asChild>
                      <a
                        href={`/problems/${item.id}`}
                        className="flex items-center justify-between gap-2"
                      >
                        <span>{item.title}</span>
                      </a>
                    </SidebarMenuButton>
                    <SidebarMenuAction showOnHover>
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Ellipsis className="w-4 h-4" />
                        </DropdownMenuTrigger>
                        <DropdownMenuContent>
                          <DropdownMenuItem
                            variant="destructive"
                            onClick={async () => {
                              try {
                                await deleteProblem(item.id);
                                toast.success("課題を削除しました");
                              } catch (_error) {
                                toast.error("削除に失敗しました");
                              } finally {
                                mutateProblems();
                              }
                            }}
                          >
                            <span>Delete</span>
                          </DropdownMenuItem>
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </SidebarMenuAction>
                  </SidebarMenuItem>
                ))}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        )}
      </SidebarContent>
      <SidebarFooter>
        <SidebarMenu>
          <SidebarMenuItem>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <SidebarMenuButton>
                  {user?.photoURL ? (
                    <Image
                      src={user?.photoURL}
                      alt="User"
                      className="w-4 h-4"
                    />
                  ) : (
                    <User2 />
                  )}
                  <span>{user?.email}</span>
                  <ChevronUp className="ml-auto" />
                </SidebarMenuButton>
              </DropdownMenuTrigger>
              <DropdownMenuContent
                side="top"
                className="w-[--radix-popper-anchor-width]"
              >
                <DropdownMenuItem
                  variant="destructive"
                  onClick={() => signOut(auth)}
                >
                  <span>ログアウト</span>
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarFooter>
    </Sidebar>
  );
}
