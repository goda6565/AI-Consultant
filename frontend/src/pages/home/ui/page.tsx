"use client";

import { ProblemForm } from "@/pages/home/ui/form";
import { Heading } from "@/shared/ui";

export function HomePage() {
  return (
    <div className="flex flex-col gap-5 h-full">
      <div className="flex gap-2 justify-between items-center">
        <Heading>Home</Heading>
      </div>
      <ProblemForm />
    </div>
  );
}
