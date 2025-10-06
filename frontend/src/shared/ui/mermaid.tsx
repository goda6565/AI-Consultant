"use client";

import mermaid from "mermaid";
import { useEffect, useRef } from "react";

type Props = { chart: string };

export function MermaidMindMap({ chart }: Props) {
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!ref.current) return;
    mermaid.initialize({ startOnLoad: true });
    mermaid.contentLoaded(); // Mermaidに再スキャンさせる
  }, []);

  return (
    <div ref={ref} className="mermaid">
      {chart}
    </div>
  );
}
