import Link from "next/link";
import { ReactNode } from "react";
import { SearchBar } from "./SearchBar";

const NAV = [
  { href: "/providers", label: "Providers" },
  { href: "/policies", label: "Policy packs" },
  { href: "/frameworks", label: "Frameworks" },
];

export function RegistryShell({
  children,
  active,
}: {
  children: ReactNode;
  active?: string;
}) {
  return (
    <div className="min-h-screen flex flex-col">
      <header className="border-b border-[var(--registry-border)] bg-white">
        <div className="max-w-7xl mx-auto px-4 py-3 flex items-center gap-6">
          <Link href="/" className="flex items-center gap-2 shrink-0">
            <span className="font-semibold text-slate-900">Faramesh</span>
            <span className="text-slate-400">|</span>
            <span className="text-slate-600">Registry</span>
          </Link>
          <div className="flex-1 max-w-xl">
            <SearchBar />
          </div>
          <nav className="hidden md:flex items-center gap-4 text-sm text-slate-600">
            <a
              href="https://docs.faramesh.dev/registry/"
              className="hover:text-[var(--registry-purple)]"
              target="_blank"
              rel="noreferrer"
            >
              Publish guide
            </a>
          </nav>
        </div>
        <div className="max-w-7xl mx-auto px-4 flex gap-1 border-t border-[var(--registry-border)]">
          {NAV.map((item) => (
            <Link
              key={item.href}
              href={item.href}
              className={`px-4 py-2.5 text-sm font-medium border-b-2 -mb-px ${
                active === item.href
                  ? "border-[var(--registry-purple)] text-[var(--registry-purple)]"
                  : "border-transparent text-slate-600 hover:text-slate-900"
              }`}
            >
              {item.label}
            </Link>
          ))}
        </div>
      </header>
      <main className="flex-1 bg-[var(--registry-surface)]">{children}</main>
      <footer className="border-t border-[var(--registry-border)] bg-white py-8 text-sm text-slate-500">
        <div className="max-w-7xl mx-auto px-4 flex flex-wrap gap-6">
          <Link href="https://docs.faramesh.dev" className="hover:text-slate-800">
            Docs
          </Link>
          <span className="ml-auto text-slate-400">
            Policy and framework artifacts are FPL-first; YAML/JSON when published as sidecars.
          </span>
        </div>
      </footer>
    </div>
  );
}
