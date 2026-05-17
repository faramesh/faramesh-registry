const API = "/api/registry";

async function registryFetch(url: string, init?: RequestInit) {
  const ctrl = new AbortController();
  const timer = setTimeout(() => ctrl.abort(), 5000);
  try {
    return await fetch(url, { ...init, signal: ctrl.signal });
  } finally {
    clearTimeout(timer);
  }
}

export type PackSummary = {
  kind: string;
  name: string;
  latest_version: string;
  description: string;
  trust_tier?: string;
  category?: string;
};

export type PackVersion = {
  kind: string;
  name: string;
  version: string;
  description?: string;
  policy_fpl: string;
  policy_yaml?: string;
  policy_json?: string;
  sha256_hex: string;
  trust_tier?: string;
  readme_markdown?: string;
  changelog?: string;
  rules_summary?: { permit?: string[]; deny?: string[]; defer?: string[] };
};

export type ProviderVersion = {
  kind: string;
  name: string;
  version: string;
  dev_only?: boolean;
  capabilities?: string[];
  downloads: Record<string, { url: string; sha256_hex: string }>;
  readme_markdown?: string;
};

export async function search(q = "", kind = ""): Promise<PackSummary[]> {
  const params = new URLSearchParams();
  if (q) params.set("q", q);
  if (kind) params.set("kind", kind);
  const res = await registryFetch(`${API}/v1/search?${params}`, { next: { revalidate: 30 } });
  if (!res.ok) throw new Error(`search failed: ${res.status}`);
  const data = await res.json();
  return data.packs ?? [];
}

function versionURL(kind: string, name: string, version: string): string {
  const base =
    kind === "provider"
      ? "providers"
      : kind === "framework"
        ? "frameworks"
        : "policies";
  return `${API}/v1/${base}/${name}/versions/${version}`;
}

export async function fetchPolicy(name: string, version: string): Promise<PackVersion> {
  const res = await registryFetch(versionURL("policy", name, version), {
    next: { revalidate: 60 },
  });
  if (!res.ok) throw new Error(`policy fetch failed: ${res.status}`);
  return res.json();
}

export async function fetchFramework(name: string, version: string): Promise<PackVersion> {
  const res = await registryFetch(versionURL("framework", name, version), {
    next: { revalidate: 60 },
  });
  if (!res.ok) throw new Error(`framework fetch failed: ${res.status}`);
  return res.json();
}

export async function fetchProvider(name: string, version: string): Promise<ProviderVersion> {
  const res = await registryFetch(versionURL("provider", name, version), {
    next: { revalidate: 60 },
  });
  if (!res.ok) throw new Error(`provider fetch failed: ${res.status}`);
  return res.json();
}

export function importSnippet(kind: string, name: string, version: string): string {
  const segment =
    kind === "provider" ? "providers" : kind === "framework" ? "frameworks" : "policies";
  return `import "registry.faramesh.dev/${segment}/${name}@${version}"`;
}

export function providerBlock(name: string): string {
  const short = name.includes("/") ? name.split("/").pop() : name;
  return `provider "${short}" {\n  type = "${short}"\n  # configure per provider docs\n}`;
}
