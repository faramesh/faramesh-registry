# GitHub distribution (before registry.faramesh.dev)

Use **one repo** — this `faramesh-registry` tree inside [Faramesh-Nexus](https://github.com/faramesh/Faramesh-Nexus). Do **not** split each provider into its own repo; the CLI resolves everything from `catalog/catalog.json`.

## What ships today

| Kind | Built? | In catalog | Notes |
|------|--------|------------|--------|
| **Provider** `faramesh/vault` | Go sidecar binary | yes | Secrets (HashiCorp Vault) |
| **Provider** `faramesh/spiffe` | Go sidecar binary | yes | Identity (SPIRE/SPIFFE) |
| **Provider** `faramesh/dev-kms` | Go sidecar binary | yes | Dev-only KMS; not for production |
| **Policy** `faramesh/demo`, `stripe`, `shell` | FPL only | yes | No binary |
| **Framework** `langgraph`, `mcp`, `cursor`, `langchain` | FPL profile only | yes | Wiring fragment, not a binary |

Future providers (AWS SM, GCP SM, audit sinks, etc.) are specified in `docs/internal/FARAMESH_REGISTRY_PLATFORM.md` — add them here with the same layout.

## Build signed artifacts (maintainers)

```bash
cd faramesh-registry
make providers                    # linux/darwin × amd64/arm64
export REGISTRY_SIGNING_KEY_B64=...  # from: go run ./cmd/gen-signing-key
make sign-all                     # FPL + provider binaries + manifest.json
./scripts/validate-catalog.sh
```

Provider binaries land under `catalog/artifacts/providers/faramesh/<name>/1.0.0/bin/`. They are **Git LFS** objects (see `.gitattributes`).

## How users consume (no live registry platform)

### Option A — Local catalog checkout (recommended for prod today)

```bash
git clone https://github.com/faramesh/Faramesh-Nexus.git
export FARAMESH_REGISTRY_ROOT="$PWD/Faramesh-Nexus/faramesh-registry"
cd your-agent-repo
faramesh registry list
faramesh check    # resolves imports from the catalog tree
faramesh apply    # installs provider binaries from catalog/bin/
```

Equivalent: `export FARAMESH_REGISTRY_URL=file://$PWD/Faramesh-Nexus/faramesh-registry`

### Option B — Self-hosted HTTP registry

```bash
cd faramesh-registry
make sign-all   # binaries must exist under catalog/.../bin/
go run ./cmd/registry -catalog catalog -listen :9876
export FARAMESH_REGISTRY_URL=http://127.0.0.1:9876
faramesh check
```

Deploy the same image to Render (`render.yaml`) for a team-wide URL.

### Option C — Official host (when live)

```bash
export FARAMESH_REGISTRY_URL=https://registry.faramesh.dev
```

Imports stay the same: `import "registry.faramesh.dev/providers/faramesh/vault@1.0.0"`

## CLI browse / search

```bash
faramesh registry url
faramesh registry list
faramesh registry list --kind provider
faramesh registry search stripe
faramesh registry info frameworks/langgraph@1.0.0
```

## Example `governance.fms`

```hcl
trust {
  key "registry.faramesh.dev" ed25519:QHApN8LymAWpwUlmsFXW0yNcC1hXtgAcwKIgJsOLnJA=
}

import "registry.faramesh.dev/frameworks/langgraph@1.0.0"
import "registry.faramesh.dev/policies/faramesh/stripe@1.0.0" as stripe_rules
import "registry.faramesh.dev/providers/faramesh/vault@1.0.0"

provider "vault" {
  type  = "vault"
  addr  = env("VAULT_ADDR")
  token = env("VAULT_TOKEN")
  mount = "secret"
}
```

Faramesh compiles FPL at `faramesh check` (policy/framework imports) and downloads provider binaries at `faramesh apply`.

## Community artifacts

1. Fork Faramesh-Nexus (or publish your own registry server).
2. Add under `faramesh-registry/catalog/artifacts/...` following existing layout.
3. Register in `catalog/catalog.json` with `trust_tier: "community"`.
4. Open a PR; CI runs `validate-catalog.sh`.
5. Consumers add **your** trust key in `trust { ... }` and either:
   - point `FARAMESH_REGISTRY_ROOT` at their clone of your fork, or
   - run your registry HTTP endpoint and set `FARAMESH_REGISTRY_URL`.

There is no central approval UI yet — trust is explicit in `governance.fms`.
