import { notFound } from "next/navigation";
import { RegistryShell } from "@/components/RegistryShell";
import { CopyBox } from "@/components/CopyBox";
import { fetchProvider, importSnippet, providerBlock } from "@/lib/registry";

type Props = { params: Promise<{ slug: string[] }> };

export default async function ProviderDetailPage({ params }: Props) {
  const { slug } = await params;
  if (!slug || slug.length < 3 || slug[slug.length - 2] !== "versions") {
    notFound();
  }
  const version = slug[slug.length - 1];
  const name = slug.slice(0, -2).join("/");
  const pv = await fetchProvider(name, version).catch(() => null);
  if (!pv) notFound();

  const imp = importSnippet("provider", name, version);

  return (
    <RegistryShell active="/providers">
      <div className="max-w-7xl mx-auto px-4 py-8">
        <nav className="text-sm text-slate-500 mb-4">
          <a href="/providers" className="hover:text-[var(--registry-purple)]">
            Providers
          </a>
          <span className="mx-2">/</span>
          <span>{name}</span>
          <span className="mx-2">/</span>
          <span>{version}</span>
        </nav>
        <h1 className="text-3xl font-bold flex items-center gap-3 flex-wrap">
          {name}
          {pv.dev_only && (
            <span className="text-xs font-semibold uppercase tracking-wide px-2 py-1 rounded bg-amber-100 text-amber-900 border border-amber-300">
              Dev only
            </span>
          )}
        </h1>
        <p className="mt-2 text-slate-600">
          Capabilities: {(pv.capabilities ?? []).join(", ") || "—"}
        </p>
        <div className="mt-8 grid lg:grid-cols-3 gap-8">
          <div className="lg:col-span-2 space-y-4">
            {pv.readme_markdown && (
              <pre className="text-sm whitespace-pre-wrap registry-card p-4 text-slate-700">
                {pv.readme_markdown}
              </pre>
            )}
            <div className="registry-card overflow-hidden">
              <table className="w-full text-sm">
                <thead className="bg-slate-50 text-slate-500 text-left">
                  <tr>
                    <th className="p-3">Platform</th>
                    <th className="p-3">SHA-256</th>
                  </tr>
                </thead>
                <tbody>
                  {Object.entries(pv.downloads).map(([plat, dl]) => (
                    <tr key={plat} className="border-t border-[var(--registry-border)]">
                      <td className="p-3 font-mono">{plat}</td>
                      <td className="p-3 font-mono text-xs text-slate-500 break-all">
                        {dl.sha256_hex}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
          <aside className="space-y-4">
            <CopyBox label="Import (FPL)" text={imp} />
            <CopyBox label="Provider block (FPL)" text={providerBlock(name)} />
          </aside>
        </div>
      </div>
    </RegistryShell>
  );
}
