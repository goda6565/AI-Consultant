"use client";

import { signOut } from "firebase/auth";
import { ChevronUp, FileText, Plus, Search, User2 } from "lucide-react";
import Image from "next/image";
import { useEffect, useState } from "react";
import { auth } from "@/shared/config";
import { useAuth } from "@/shared/hooks";
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
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarTrigger,
  useSidebar,
} from "@/shared/ui";

const menuItems = [
  {
    title: "新しいケース",
    url: "/",
    icon: Plus,
  },
  {
    title: "ケースを検索",
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
  const user = useAuth();

  // サイドバーが開かれた時にホバー状態をリセット
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
      </SidebarContent>
      <SidebarFooter>
        <SidebarMenu>
          <SidebarMenuItem>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <SidebarMenuButton>
                  {user.user?.photoURL ? (
                    <Image
                      src={user.user?.photoURL}
                      alt="User"
                      className="w-4 h-4"
                    />
                  ) : (
                    <User2 />
                  )}
                  <span>{user.user?.email}</span>
                  <ChevronUp className="ml-auto" />
                </SidebarMenuButton>
              </DropdownMenuTrigger>
              <DropdownMenuContent
                side="top"
                className="w-[--radix-popper-anchor-width]"
              >
                <DropdownMenuItem onClick={() => signOut(auth)}>
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
