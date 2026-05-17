# Catalog production status (honest)

## What works today (verified)

| Flow | Status |
|------|--------|
| `faramesh registry list` / `search` | Fetches `catalog/catalog.json` from GitHub |
| `faramesh check` with policy/framework imports | Downloads FPL from GitHub, merges AST |
| `faramesh apply` provider import | Downloads platform binary from GitHub raw, sha256 check |
| Provider binaries in repo | `vault`, `spiffe`, `dev-kms` — four OS/arch builds each |

## Artifact shapes (this is normal)

| Kind | Files | Why it looks “small” |
|------|-------|----------------------|
| **Framework profile** | `profile.fpl` (+ optional `meta.json`) | Wiring fragment only — not a full stack |
| **Policy pack** | `policy.fpl` (+ `.sig`, `meta.json`) | Real rules live in `policy.fpl`; JSON is metadata |
| **Provider** | `manifest.json` + `bin/*` (+ `.sig`) | Binary is the artifact; manifest is the index |

Example framework profile (`langgraph@1.0.0`) is ~10 lines: sets `framework "langgraph"` and default posture. That is intentional.

## Production readiness by artifact

| Artifact | Production? | Notes |
|----------|-------------|-------|
| `faramesh/vault` | **Yes** (with your Vault) | Sidecar + manifest hashes; needs live Vault at apply |
| `faramesh/spiffe` | **Yes** (with SPIRE) | Same |
| `faramesh/dev-kms` | **No** | Dev/test only — marked in catalog |
| `faramesh/stripe`, `shell` | **Starter packs** | Real FPL; tune thresholds for your org |
| `faramesh/openai`, `github`, `mcp` | **Starter baselines** | Review rules before prod; not compliance-certified |
| Framework profiles | **Yes as imports** | Pin version; extend in your own `agent` blocks |

## Signatures

- Official **public** key: `catalog/trust/keys.json` (`faramesh-ed25519-2026`).
- Publishing that public key in docs is correct — only `REGISTRY_SIGNING_KEY_B64` (private) must stay in GitHub Actions secrets.
- CI signing runs only when that secret is set. See [SETUP_SIGNING.md](./SETUP_SIGNING.md).

## Browse UI

The Next.js app under `web/` is **not** the supported distribution path. Use GitHub + `faramesh registry` CLI. The UI is disabled in deploy configs.
