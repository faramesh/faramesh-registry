import { Suspense } from "react";
import { RegistryShell } from "@/components/RegistryShell";
import { BrowseSidebar } from "@/components/BrowseSidebar";
import { CatalogCard } from "@/components/CatalogCard";
import { search } from "@/lib/registry";

export default async function PoliciesPage({
  searchParams,
}: {
  searchParams: Promise<{ tier?: string; q?: string }>;
}) {
  const sp = await searchParams;
  const items = await search(sp.q ?? "", "policy").catch(() => []);
  const filtered = sp.tier
    ? items.filter((i) => (i.trust_tier ?? "").toLowerCase() === sp.tier)
    : items;

  return (
    <RegistryShell active="/policies">
      <div className="max-w-7xl mx-auto px-4 py-8">
        <h1 className="text-2xl font-bold text-slate-900">Policy packs</h1>
        <p className="mt-2 text-slate-600 max-w-3xl">
          Versioned FPL rule sets. The registry stores <code className="text-sm bg-slate-100 px-1">policy.fpl</code> as
          canonical bytes; optional YAML/JSON sidecars appear on the artifact page when published.
        </p>
        <div className="mt-8 flex gap-8">
          <Suspense fallback={null}>
            <BrowseSidebar kind="policies" />
          </Suspense>
          <div className="flex-1">
            <p className="text-sm text-slate-500 mb-4">
              Showing {filtered.length} policy pack{filtered.length === 1 ? "" : "s"}
            </p>
            <div className="grid sm:grid-cols-2 gap-4">
              {filtered.map((item) => (
                <CatalogCard key={item.name} item={item} />
              ))}
            </div>
            {filtered.length === 0 && (
              <p className="text-slate-500 text-sm">No policy packs match your filters.</p>
            )}
          </div>
        </div>
      </div>
    </RegistryShell>
  );
}
