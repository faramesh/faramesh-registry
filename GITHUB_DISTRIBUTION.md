# GitHub catalog distribution

All Faramesh registry artifacts live in this repository. The Faramesh CLI resolves imports from GitHub by default.

## What ships in the catalog

| Kind | Examples | Binary? |
|------|----------|---------|
| **Providers** | `faramesh/vault`, `faramesh/spiffe` | Yes — gRPC sidecar |
| **Policy packs** | `faramesh/stripe`, `faramesh/openai`, `faramesh/github` | FPL only |
| **Framework profiles** | `langgraph`, `mcp`, `cursor`, `crewai`, `bedrock` | FPL only |

`faramesh/dev-kms` is included for local development only — do not use in production.

## Default CLI behavior

With no environment variables set, `faramesh check` and `faramesh apply` fetch:

- `catalog/catalog.json` from `raw.githubusercontent.com`
- FPL files and provider binaries from the same repo/ref (`main` by default)

## Import syntax

```hcl
trust {
  key "github.com/faramesh/faramesh-registry" ed25519:QHApN8LymAWpwUlmsFXW0yNcC1hXtgAcwKIgJsOLnJA=
}

import "github.com/faramesh/faramesh-registry/frameworks/langgraph@1.0.0"
import "github.com/faramesh/faramesh-registry/policies/faramesh/stripe@1.0.0" as stripe_rules
import "github.com/faramesh/faramesh-registry/providers/faramesh/vault@1.0.0"
```

The `ed25519:...` value above is the **public** verification key (same as `catalog/trust/keys.json`). It is safe to publish. The **private** signing key lives only in GitHub Actions as `REGISTRY_SIGNING_KEY_B64` — see [SETUP_SIGNING.md](./SETUP_SIGNING.md).

## Overrides

| Mode | Setup |
|------|--------|
| **GitHub (default)** | Nothing to configure |
| **Local clone** | `export FARAMESH_REGISTRY_ROOT=$PWD` after cloning this repo |
| **HTTP server** | `go run ./cmd/registry -catalog catalog` and `export FARAMESH_REGISTRY_URL=http://127.0.0.1:9876` |
| **Fork / community** | Clone fork, add trust key, set `FARAMESH_REGISTRY_ROOT` |

## Maintainer workflow

```bash
make providers
./scripts/refresh-provider-hashes.sh
./scripts/validate-catalog.sh
# With signing key:
export REGISTRY_SIGNING_KEY_B64=...
make sign-all
git push
```

## Community publishing

1. Fork this repository.
2. Add artifacts under `catalog/artifacts/` and register in `catalog/catalog.json`.
3. Open a pull request; CI validates FPL and manifest layout.
4. Consumers pin your fork via `FARAMESH_REGISTRY_ROOT` and add your publisher key to `trust { ... }`.
