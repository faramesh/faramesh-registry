import Link from "next/link";
import type { PackSummary } from "@/lib/registry";

export function ArtifactCard({ item }: { item: PackSummary }) {
  const href = `/${item.kind}s/${item.name}/versions/${item.latest_version}`;
  return (
    <Link
      href={href}
      className="block rounded-xl border border-zinc-800 bg-zinc-900/50 p-5 hover:border-teal-800 transition-colors"
    >
      <div className="flex items-start justify-between gap-2">
        <span className="text-xs uppercase text-zinc-500">{item.kind}</span>
        {item.trust_tier && (
          <span className="text-[10px] uppercase px-2 py-0.5 rounded bg-teal-950 text-teal-400 border border-teal-900">
            {item.trust_tier}
          </span>
        )}
      </div>
      <h3 className="mt-2 font-semibold text-lg">{item.name}</h3>
      <p className="mt-1 text-sm text-zinc-400 line-clamp-2">{item.description}</p>
      <p className="mt-3 text-xs text-zinc-500">latest {item.latest_version}</p>
    </Link>
  );
}
