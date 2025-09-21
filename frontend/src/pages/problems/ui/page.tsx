"use client";

import { redirect } from "next/navigation";
import { use } from "react";
import { useGetProblem } from "@/shared/api";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
  Badge,
  Heading,
  LoadingPage,
} from "@/shared/ui";
import { Chat } from "./chat";

type ProblemPageProps = {
  params: Promise<{ id: string }>;
};

export function ProblemPage({ params }: ProblemPageProps) {
  const { id } = use(params);
  const {
    data: problem,
    isLoading: isProblemLoading,
    mutate: mutateProblem,
  } = useGetProblem(id);
  if (isProblemLoading) {
    return <LoadingPage />;
  }
  if (!problem) {
    redirect("/");
  }
  return (
    <div className="flex flex-col h-full">
      <div className="flex gap-6 justify-between items-center p-4 border-b flex-shrink-0">
        <Accordion type="single" collapsible className="w-full">
          <AccordionItem value="item-1">
            <AccordionTrigger>
              <Heading>{problem.title}</Heading>
            </AccordionTrigger>
            <AccordionContent>{problem.description}</AccordionContent>
          </AccordionItem>
        </Accordion>
        <Badge variant="outline">{problem.status}</Badge>
      </div>
      <div className="flex-1 min-h-0">
        <Chat problem={problem} mutateProblem={mutateProblem} />
      </div>
    </div>
  );
}
