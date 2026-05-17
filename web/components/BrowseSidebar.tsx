"use client";

import { useRouter, useSearchParams } from "next/navigation";

const TIERS = ["official", "partner", "community"];

export function BrowseSidebar({ kind }: { kind: string }) {
  const router = useRouter();
  const params = useSearchParams();
  const tier = params.get("tier") ?? "";

  function setTier(t: string) {
    const next = new URLSearchParams(params.toString());
    if (t) next.set("tier", t);
    else next.delete("tier");
    router.push(`/${kind}?${next.toString()}`);
  }

  return (
    <aside className="w-56 shrink-0 space-y-6">
      <div>
        <h3 className="text-xs font-semibold uppercase tracking-wide text-slate-500 mb-2">
          Tier
        </h3>
        <ul className="space-y-1 text-sm">
          <li>
            <button
              type="button"
              onClick={() => setTier("")}
              className={`${!tier ? "text-[var(--registry-purple)] font-medium" : "text-slate-600 hover:text-slate-900"}`}
            >
              All
            </button>
          </li>
          {TIERS.map((t) => (
            <li key={t}>
              <button
                type="button"
                onClick={() => setTier(t)}
                className={`capitalize ${tier === t ? "text-[var(--registry-purple)] font-medium" : "text-slate-600 hover:text-slate-900"}`}
              >
                {t}
              </button>
            </li>
          ))}
        </ul>
      </div>
    </aside>
  );
}
