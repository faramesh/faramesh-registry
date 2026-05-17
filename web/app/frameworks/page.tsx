import { Suspense } from "react";
import { RegistryShell } from "@/components/RegistryShell";
import { BrowseSidebar } from "@/components/BrowseSidebar";
import { CatalogCard } from "@/components/CatalogCard";
import { search } from "@/lib/registry";

export default async function FrameworksPage({
  searchParams,
}: {
  searchParams: Promise<{ tier?: string }>;
}) {
  const sp = await searchParams;
  const items = await search("", "framework").catch(() => []);
  const filtered = sp.tier
    ? items.filter((i) => (i.trust_tier ?? "").toLowerCase() === sp.tier)
    : items;

  return (
    <RegistryShell active="/frameworks">
      <div className="max-w-7xl mx-auto px-4 py-8">
        <h1 className="text-2xl font-bold text-slate-900">Framework profiles</h1>
        <p className="mt-2 text-slate-600 max-w-3xl">
          FPL fragments (<code className="text-sm bg-slate-100 px-1">profile.fpl</code>) that wire SDK/MCP
          enforcement for agent frameworks. Import in <code className="text-sm bg-slate-100 px-1">governance.fms</code>{" "}
          and pair with SDK autopatch.
        </p>
        <div className="mt-8 flex gap-8">
          <Suspense fallback={null}>
            <BrowseSidebar kind="frameworks" />
          </Suspense>
          <div className="flex-1">
            <div className="grid sm:grid-cols-2 gap-4">
              {filtered.map((item) => (
                <CatalogCard key={item.name} item={item} />
              ))}
            </div>
            {filtered.length === 0 && (
              <p className="text-slate-500 text-sm">No framework profiles match your filters.</p>
            )}
          </div>
        </div>
      </div>
    </RegistryShell>
  );
}
