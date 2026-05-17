import Link from "next/link";
import { RegistryShell } from "@/components/RegistryShell";
import { CatalogCard } from "@/components/CatalogCard";
import { search } from "@/lib/registry";

export default async function HomePage() {
  const all = await search().catch(() => []);
  const policies = all.filter((p) => p.kind === "policy");
  const frameworks = all.filter((p) => p.kind === "framework");

  return (
    <RegistryShell>
      <div className="max-w-7xl mx-auto px-4 py-12">
        <section className="text-center max-w-2xl mx-auto">
          <h1 className="text-4xl font-bold text-slate-900 tracking-tight">
            Faramesh Registry
          </h1>
          <p className="mt-4 text-slate-600">
            Discover providers, policy packs, and framework profiles for{" "}
            <code className="text-sm bg-slate-100 px-1 rounded">governance.fms</code>.
            Artifacts are authored in <strong>FPL</strong> first; YAML and JSON sidecars when
            publishers include them.
          </p>
          <div className="mt-8 flex flex-wrap justify-center gap-3">
            <Link
              href="/policies"
              className="px-5 py-2.5 rounded-md bg-[var(--registry-purple)] text-white text-sm font-medium hover:opacity-90"
            >
              Browse policy packs
            </Link>
            <Link
              href="/frameworks"
              className="px-5 py-2.5 rounded-md border border-[var(--registry-border)] text-sm font-medium hover:bg-white"
            >
              Browse frameworks
            </Link>
            <Link
              href="/providers"
              className="px-5 py-2.5 rounded-md border border-[var(--registry-border)] text-sm font-medium hover:bg-white"
            >
              Providers
            </Link>
          </div>
          <p className="mt-6 text-sm text-slate-500">
            {policies.length} policy packs · {frameworks.length} framework profiles in catalog
          </p>
        </section>

        {policies.length > 0 && (
          <section className="mt-16">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-semibold">Featured policy packs</h2>
              <Link href="/policies" className="text-sm text-[var(--registry-purple)]">
                See all →
              </Link>
            </div>
            <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-4">
              {policies.slice(0, 6).map((item) => (
                <CatalogCard key={item.name} item={item} />
              ))}
            </div>
          </section>
        )}
      </div>
    </RegistryShell>
  );
}
