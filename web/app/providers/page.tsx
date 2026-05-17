import { Suspense } from "react";
import { RegistryShell } from "@/components/RegistryShell";
import { BrowseSidebar } from "@/components/BrowseSidebar";
import { CatalogCard } from "@/components/CatalogCard";
import { search } from "@/lib/registry";

export default async function ProvidersPage({
  searchParams,
}: {
  searchParams: Promise<{ tier?: string }>;
}) {
  const sp = await searchParams;
  const items = await search("", "provider").catch(() => []);
  const filtered = sp.tier
    ? items.filter((i) => (i.trust_tier ?? "").toLowerCase() === sp.tier)
    : items;

  return (
    <RegistryShell active="/providers">
      <div className="max-w-7xl mx-auto px-4 py-8">
        <h1 className="text-2xl font-bold text-slate-900">Providers</h1>
        <p className="mt-2 text-slate-600 max-w-3xl">
          Signed <code className="text-sm bg-slate-100 px-1">ProviderService</code> binaries.
          Configure in FPL with an <code className="text-sm bg-slate-100 px-1">import</code> and{" "}
          <code className="text-sm bg-slate-100 px-1">provider</code> block.
        </p>
        <div className="mt-4 rounded-md border border-amber-200 bg-amber-50 p-4 text-sm text-amber-900 max-w-3xl">
          <strong>No provider binaries are published in this catalog yet.</strong> For local development,
          use built-in providers via <code className="bg-amber-100 px-1">faramesh dev</code> (for example{" "}
          <code className="bg-amber-100 px-1">type = &quot;dev-kms&quot;</code> in{" "}
          <code className="bg-amber-100 px-1">governance.fms</code>). Registry provider entries will
          appear here only after real signed binaries are uploaded through GitOps.
        </div>
        <div className="mt-8 flex gap-8">
          <Suspense fallback={null}>
            <BrowseSidebar kind="providers" />
          </Suspense>
          <div className="flex-1">
            {filtered.length > 0 ? (
              <div className="grid sm:grid-cols-2 gap-4">
                {filtered.map((item) => (
                  <CatalogCard key={item.name} item={item} />
                ))}
              </div>
            ) : (
              <p className="text-slate-500 text-sm">No providers in the catalog.</p>
            )}
          </div>
        </div>
      </div>
    </RegistryShell>
  );
}
