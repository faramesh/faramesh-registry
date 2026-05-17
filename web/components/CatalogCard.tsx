import Link from "next/link";
import type { PackSummary } from "@/lib/registry";

export function CatalogCard({ item }: { item: PackSummary }) {
  const href = `/${item.kind}s/${item.name}/versions/${item.latest_version}`;
  return (
    <Link href={href} className="registry-card block p-4 h-full">
      <div className="flex justify-between items-start gap-2">
        <span className="text-xs text-slate-500 capitalize">{item.kind}</span>
        {item.trust_tier === "official" && (
          <span className="registry-badge-official">Official</span>
        )}
      </div>
      <h3 className="mt-2 font-semibold text-slate-900">{item.name}</h3>
      <p className="mt-1 text-sm text-slate-600 line-clamp-2">{item.description}</p>
      <p className="mt-3 text-xs text-slate-400">Latest {item.latest_version}</p>
    </Link>
  );
}
