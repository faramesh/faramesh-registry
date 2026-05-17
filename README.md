# Faramesh Registry

**Website:** https://registry.faramesh.dev (production)  
**Documentation:** https://docs.faramesh.dev/registry/  
**Runtime spec:** [`faramesh-core/docs/internal/FARAMESH.md`](../faramesh-core/docs/internal/FARAMESH.md)  
**Platform design:** [`docs/internal/FARAMESH_REGISTRY_PLATFORM.md`](../docs/internal/FARAMESH_REGISTRY_PLATFORM.md)  
**Language principles:** [PRINCIPLES.md](./PRINCIPLES.md) — **FPL is primary**; YAML/JSON are natural sidecars.

## What this repo is

Official distribution for **providers** (signed binaries), **policy packs** (FPL), and **framework profiles** (FPL). The CLI resolves pinned imports at `faramesh check` / `apply`.

## GitHub distribution (no registry.faramesh.dev yet)

See **[GITHUB_DISTRIBUTION.md](./GITHUB_DISTRIBUTION.md)** — one repo, built provider binaries under `catalog/artifacts/.../bin/`, consume via `FARAMESH_REGISTRY_ROOT` or self-hosted HTTP.

```bash
export FARAMESH_REGISTRY_ROOT=/path/to/Faramesh-Nexus/faramesh-registry
faramesh registry list
faramesh check
```

## Quick start

```bash
go run ./cmd/registry -catalog catalog
# → http://127.0.0.1:9876

export FARAMESH_REGISTRY_URL=http://127.0.0.1:9876
curl -s http://127.0.0.1:9876/.well-known/faramesh.json | jq .
curl -s "http://127.0.0.1:9876/v1/policies/faramesh/demo/versions/0.1.0" | jq .policy_fpl
```

### Web UI

```bash
cd web && npm install && npm run dev
# → http://localhost:3001
```

### Docker

```bash
docker compose up --build
```

## Catalog (GitOps)

- Index: `catalog/catalog.json`
- Canonical policy/framework files: `policy.fpl`, `profile.fpl`
- Optional: `policy.yaml`, `policy.json`, `meta.json`, `README.md`
- Validate: `./scripts/validate-catalog.sh`
- Sign FPL: `go run ./cmd/sign-catalog -catalog catalog` (requires `REGISTRY_SIGNING_KEY_B64`)

## API (v1)

| Method | Path |
|--------|------|
| GET | `/.well-known/faramesh.json` |
| GET | `/v1/search?q=&kind=&tier=` |
| GET | `/v1/stats` |
| GET | `/v1/policies/{name}/versions/{version}` — JSON includes `policy_fpl`; `?format=yaml\|json` for sidecars |
| GET | `/v1/frameworks/{name}/versions/{version}` |
| GET | `/v1/providers/{name}/versions/{version}` |
| GET | `/v1/trust/keys` |
| POST | `/v1/publish` — verifies Ed25519 over `policy_fpl` |

Namespaced artifacts use slash paths, e.g. `faramesh/stripe`.

## Implementation status

See [IMPLEMENTATION_PLAN.md](./IMPLEMENTATION_PLAN.md) — R0–R5 done; R6 (community tier) planned.
