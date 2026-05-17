import { notFound } from "next/navigation";
import { RegistryShell } from "@/components/RegistryShell";
import { CopyBox } from "@/components/CopyBox";
import { FormatViewer } from "@/components/FormatViewer";
import { fetchPolicy, importSnippet } from "@/lib/registry";

type Props = { params: Promise<{ slug: string[] }> };

export default async function PolicyDetailPage({ params }: Props) {
  const { slug } = await params;
  if (!slug || slug.length < 3 || slug[slug.length - 2] !== "versions") {
    notFound();
  }
  const version = slug[slug.length - 1];
  const name = slug.slice(0, -2).join("/");
  const pack = await fetchPolicy(name, version).catch(() => null);
  if (!pack) notFound();

  const imp = importSnippet("policy", name, version);
  return (
    <RegistryShell active="/policies">
      <div className="max-w-7xl mx-auto px-4 py-8">
        <nav className="text-sm text-slate-500 mb-4">
          <a href="/policies" className="hover:text-[var(--registry-purple)]">
            Policy packs
          </a>
          <span className="mx-2">/</span>
          <span className="text-slate-800">{name}</span>
          <span className="mx-2">/</span>
          <span>{version}</span>
        </nav>

        <div className="flex flex-wrap items-start justify-between gap-4 border-b border-[var(--registry-border)] pb-6">
          <div>
            <div className="flex items-center gap-2">
              <h1 className="text-3xl font-bold text-slate-900">{name}</h1>
              {pack.trust_tier === "official" && (
                <span className="registry-badge-official">Official</span>
              )}
            </div>
            <p className="mt-2 text-slate-600">{pack.description}</p>
            <p className="mt-2 text-xs text-slate-400 font-mono">
              sha256:{pack.sha256_hex.slice(0, 16)}…
            </p>
          </div>
          <div className="text-sm text-slate-600">
            Version <strong>{version}</strong>
          </div>
        </div>

        <div className="mt-8 grid lg:grid-cols-3 gap-8">
          <div className="lg:col-span-2 space-y-6">
            <section>
              <h2 className="text-lg font-semibold mb-3">Policy source</h2>
              <FormatViewer
                fpl={pack.policy_fpl}
                yaml={pack.policy_yaml}
                json={pack.policy_json}
              />
            </section>
            {pack.rules_summary && (
              <section className="registry-card p-4">
                <h3 className="font-medium text-slate-900">Rules summary</h3>
                <ul className="mt-2 text-sm text-slate-600 space-y-1">
                  {pack.rules_summary.permit?.map((r) => (
                    <li key={`p-${r}`}>permit {r}</li>
                  ))}
                  {pack.rules_summary.deny?.map((r) => (
                    <li key={`d-${r}`}>deny {r}</li>
                  ))}
                </ul>
              </section>
            )}
          </div>
          <aside className="space-y-4">
            <CopyBox label="Import (FPL)" text={imp} />
            <CopyBox
              label="Example stack snippet"
              text={`${imp} as rules\n\nagent "my-agent" {\n  use rules {}\n}`}
            />
            <div className="registry-card p-4 text-sm text-slate-600">
              <h3 className="font-medium text-slate-900 mb-2">How to use</h3>
              <ol className="list-decimal list-inside space-y-1">
                <li>Add the import to <code className="text-xs bg-slate-100 px-1">governance.fms</code></li>
                <li>Run <code className="text-xs bg-slate-100 px-1">faramesh check</code></li>
                <li>Run <code className="text-xs bg-slate-100 px-1">faramesh apply</code></li>
              </ol>
            </div>
          </aside>
        </div>
      </div>
    </RegistryShell>
  );
}
