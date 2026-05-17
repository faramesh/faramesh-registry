"use client";

export function CopyBox({ label, text }: { label: string; text: string }) {
  return (
    <div className="rounded-lg border border-zinc-800 bg-zinc-900 p-4">
      <div className="flex items-center justify-between mb-2">
        <span className="text-xs font-medium uppercase tracking-wide text-zinc-500">
          {label}
        </span>
        <button
          type="button"
          className="text-xs text-teal-400 hover:text-teal-300"
          onClick={() => navigator.clipboard.writeText(text)}
        >
          Copy
        </button>
      </div>
      <pre className="text-sm font-mono text-zinc-200 whitespace-pre-wrap">{text}</pre>
    </div>
  );
}
