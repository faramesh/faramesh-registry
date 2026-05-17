# Faramesh Registry

Official **GitHub catalog** for Faramesh artifacts: signed provider binaries, policy packs (FPL), and framework profiles (FPL).

**Catalog:** [github.com/faramesh/faramesh-registry](https://github.com/faramesh/faramesh-registry)  
**Documentation:** [docs.faramesh.dev/registry](https://docs.faramesh.dev/registry/)  
**Runtime spec:** [faramesh-core](https://github.com/faramesh/faramesh-core)

## Use artifacts in your stack

Imports in `governance.fms` point at this repository (pinned semver required):

```hcl
import "github.com/faramesh/faramesh-registry/frameworks/langgraph@1.0.0"
import "github.com/faramesh/faramesh-registry/policies/faramesh/stripe@1.0.0" as stripe_rules
import "github.com/faramesh/faramesh-registry/providers/faramesh/vault@1.0.0"
```

The Faramesh CLI resolves these at `faramesh check` (FPL merge) and `faramesh apply` (provider binaries). No separate registry service is required.

```bash
faramesh registry list
faramesh registry search vault
faramesh check
faramesh apply
```

Optional overrides:

| Variable | Purpose |
|----------|---------|
| `FARAMESH_REGISTRY_ROOT` | Path to a local clone of this repo (air-gapped / fork) |
| `FARAMESH_REGISTRY_URL` | Self-hosted HTTP registry (`go run ./cmd/registry`) |
| `FARAMESH_REGISTRY_GITHUB_REF` | Git ref (default `main`) |

See **[GITHUB_DISTRIBUTION.md](./GITHUB_DISTRIBUTION.md)** and **[CONTRIBUTING.md](./CONTRIBUTING.md)**.

## Catalog layout

```
catalog/
  catalog.json          # index
  trust/keys.json       # Ed25519 public keys
  artifacts/
    providers/.../bin/  # signed sidecar binaries
    policies/.../policy.fpl
    frameworks/.../profile.fpl
```

Validate: `./scripts/validate-catalog.sh`  
Build providers: `make providers`  
Refresh hashes: `./scripts/refresh-provider-hashes.sh`  
Sign (maintainers): `make sign-all` with `REGISTRY_SIGNING_KEY_B64`

## Self-hosted HTTP registry (optional)

```bash
go run ./cmd/registry -catalog catalog -listen :9876
export FARAMESH_REGISTRY_URL=http://127.0.0.1:9876
```

Docker: `docker compose up --build`

## Production status

See **[PRODUCTION_STATUS.md](./PRODUCTION_STATUS.md)** (what is verified vs starter packs) and **[SETUP_SIGNING.md](./SETUP_SIGNING.md)** (CI signing secret).

The browse UI under `web/` is **disabled**; distribution is GitHub + CLI only.
