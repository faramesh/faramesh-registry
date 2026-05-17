"use client";

import { useState } from "react";

type Format = "fpl" | "yaml" | "json";

export function FormatViewer({
  fpl,
  yaml,
  json,
}: {
  fpl: string;
  yaml?: string;
  json?: string;
}) {
  const [format, setFormat] = useState<Format>("fpl");
  const tabs: { id: Format; label: string; disabled?: boolean; primary?: boolean }[] = [
    { id: "fpl", label: "FPL", primary: true },
    { id: "yaml", label: "YAML", disabled: !yaml },
    { id: "json", label: "JSON", disabled: !json },
  ];
  const body =
    format === "fpl" ? fpl : format === "yaml" ? yaml ?? "" : json ?? "";

  return (
    <div className="rounded-lg border border-zinc-800 bg-zinc-900/80 overflow-hidden">
      <div className="flex border-b border-zinc-800 bg-zinc-900">
        {tabs.map((t) => (
          <button
            key={t.id}
            type="button"
            disabled={t.disabled}
            onClick={() => setFormat(t.id)}
            className={`px-4 py-2 text-sm font-medium ${
              format === t.id
                ? "text-teal-400 border-b-2 border-teal-500"
                : "text-zinc-400 hover:text-zinc-200 disabled:opacity-40"
            }`}
          >
            {t.label}
            {t.primary && (
              <span className="ml-1 text-[10px] uppercase text-teal-500">primary</span>
            )}
          </button>
        ))}
      </div>
      <pre className="p-4 text-sm font-mono overflow-x-auto max-h-[32rem] text-zinc-200 whitespace-pre-wrap">
        {body || "(empty)"}
      </pre>
    </div>
  );
}
