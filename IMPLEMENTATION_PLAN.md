# Faramesh Registry — Implementation Plan

Companion spec: [`../docs/internal/FARAMESH_REGISTRY_PLATFORM.md`](../docs/internal/FARAMESH_REGISTRY_PLATFORM.md)  
Language principles: [PRINCIPLES.md](./PRINCIPLES.md) (**FPL primary**, YAML/JSON sidecars)

## Phases

| Phase | Deliverable | Status |
|-------|-------------|--------|
| **R0** | Go registry service: `.well-known`, `/v1/search`, kind routes, static `/artifacts/` | **Done** |
| **R1** | GitOps catalog + `validate-catalog.sh` + CI workflow | **Done** |
| **R2** | `POST /v1/publish` signature verify + `/v1/trust/keys` + `cmd/sign-catalog` | **Done** |
| **R3** | Official catalog seed (demo, stripe, shell, vault, dev-kms, langgraph) | **Done** |
| **R4** | Next.js browse UI (FPL-primary viewer, YAML/JSON tabs) | **Done** |
| **R5** | Production deploy (Dockerfile, docker-compose) | **Done** (local); S3/Postgres CDN TBD |
| **R6** | Partner/community tiers + report queue | Planned |

## Run locally

```bash
# API
go run ./cmd/registry -catalog catalog

# Web (separate terminal)
cd web && npm install && npm run dev

# Or both
docker compose up --build
```

```bash
export FARAMESH_REGISTRY_URL=http://127.0.0.1:9876
cd ../faramesh-core && faramesh check
```

## Sign catalog (CI / maintainers)

```bash
export REGISTRY_SIGNING_KEY_B64="<64-byte-ed25519-private-key-base64>"
go run ./cmd/sign-catalog -catalog catalog
```

## Repo layout

```
faramesh-registry/
  cmd/registry/           # HTTP server
  cmd/sign-catalog/       # Sign policy.fpl / profile.fpl
  internal/server/        # v1 API
  internal/catalog/       # GitOps loader (FPL + sidecars)
  internal/sign/          # Ed25519 verify
  internal/trust/         # Publisher keys
  catalog/                # catalog.json + artifacts/
  web/                    # Next.js UI
  scripts/validate-catalog.sh
  docker-compose.yml
```
