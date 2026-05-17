import { RegistryShell } from "@/components/RegistryShell";
import { CatalogCard } from "@/components/CatalogCard";
import { search } from "@/lib/registry";

export default async function SearchPage({
  searchParams,
}: {
  searchParams: Promise<{ q?: string; kind?: string }>;
}) {
  const sp = await searchParams;
  const q = sp.q ?? "";
  const items = await search(q, sp.kind ?? "").catch(() => []);

  return (
    <RegistryShell>
      <div className="max-w-7xl mx-auto px-4 py-8">
        <h1 className="text-2xl font-bold">Search</h1>
        <p className="text-slate-600 mt-1">
          {q ? (
            <>
              Results for <strong>&quot;{q}&quot;</strong>
            </>
          ) : (
            "Enter a query in the header search box."
          )}
        </p>
        <div className="mt-6 grid sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {items.map((item) => (
            <CatalogCard key={`${item.kind}-${item.name}`} item={item} />
          ))}
        </div>
        {q && items.length === 0 && (
          <p className="mt-4 text-slate-500 text-sm">No artifacts matched.</p>
        )}
      </div>
    </RegistryShell>
  );
}
