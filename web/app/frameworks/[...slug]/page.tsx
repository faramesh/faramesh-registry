import { notFound } from "next/navigation";
import { RegistryShell } from "@/components/RegistryShell";
import { CopyBox } from "@/components/CopyBox";
import { FormatViewer } from "@/components/FormatViewer";
import { fetchFramework, importSnippet } from "@/lib/registry";

type Props = { params: Promise<{ slug: string[] }> };

export default async function FrameworkDetailPage({ params }: Props) {
  const { slug } = await params;
  if (!slug || slug.length < 3 || slug[slug.length - 2] !== "versions") {
    notFound();
  }
  const version = slug[slug.length - 1];
  const name = slug.slice(0, -2).join("/");
  const pack = await fetchFramework(name, version).catch(() => null);
  if (!pack) notFound();

  const imp = importSnippet("framework", name, version);

  return (
    <RegistryShell active="/frameworks">
      <div className="max-w-7xl mx-auto px-4 py-8">
        <nav className="text-sm text-slate-500 mb-4">
          <a href="/frameworks" className="hover:text-[var(--registry-purple)]">
            Frameworks
          </a>
          <span className="mx-2">/</span>
          <span>{name}</span>
          <span className="mx-2">/</span>
          <span>{version}</span>
        </nav>
        <h1 className="text-3xl font-bold">{name}</h1>
        <p className="mt-2 text-slate-600">{pack.description}</p>
        <div className="mt-8 grid lg:grid-cols-3 gap-8">
          <div className="lg:col-span-2">
            <FormatViewer fpl={pack.policy_fpl} yaml={pack.policy_yaml} json={pack.policy_json} />
          </div>
          <aside className="space-y-4">
            <CopyBox label="Import (FPL)" text={imp} />
            <CopyBox
              label="Runtime (FPL)"
              text={`${imp}\n\nruntime {\n  mode    = "enforce"\n  wal_dir = "./faramesh-wal"\n}`}
            />
          </aside>
        </div>
      </div>
    </RegistryShell>
  );
}
