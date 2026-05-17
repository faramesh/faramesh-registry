# Deploying Faramesh Registry

## What you need to provide

| Item | Required for | Notes |
|------|----------------|-------|
| **Domain + TLS** | Production | e.g. `registry.faramesh.dev` with HTTPS certificate |
| **`REGISTRY_SIGNING_KEY_B64`** | Signed FPL + provider binaries | Generate with `go run ./cmd/gen-signing-key`. Public key in `catalog/trust/keys.json` (`faramesh-ed25519-2026`) |
| **Object storage** (S3/GCS) | Large provider binaries | Not required until real provider binaries are published |
| **Postgres** (optional) | Download analytics, search index | R5 uses GitOps `catalog.json` only |

## What ships today (honest scope)

- **Policies & frameworks:** Real FPL artifacts in `catalog/artifacts/` — ready to serve.
- **Providers:** `faramesh/vault`, `faramesh/spiffe`, and `faramesh/dev-kms` sidecars build with `make providers` and sign with `make sign-catalog` (requires `REGISTRY_SIGNING_KEY_B64`). Binaries live under `catalog/artifacts/providers/.../bin/` (gitignored; upload to object storage for production CDN).
- **Web UI:** Browse/search/detail for catalog content; Terraform Registry–style layout.

## Docker (single host)

```bash
cd faramesh-registry
docker compose up --build -d
```

- API: port **9876**
- Web: port **3001** (proxies `/api/registry/*` → API)

Set in `docker-compose.yml` or environment:

```bash
REGISTRY_API_URL=http://registry:9876
```

## Manual production layout

1. Build API: `go build -o bin/registry ./cmd/registry`
2. Run: `./bin/registry -catalog /var/faramesh/catalog -listen 0.0.0.0:9876`
3. Build web: `cd web && npm ci && npm run build && npm run start`
4. Put reverse proxy (nginx/Caddy) in front:
   - `/` → Next.js `:3001`
   - Or serve API at `registry.example.com` and web at `registry.example.com` with path split

## GitOps publish flow

1. Add FPL under `catalog/artifacts/policies/.../policy.fpl`
2. Update `catalog/catalog.json`
3. `./scripts/validate-catalog.sh`
4. `go run ./cmd/sign-catalog -catalog catalog` (with signing key)
5. Deploy catalog directory to server (or CI sync to S3 + index)

## Health checks

```bash
curl -sf http://127.0.0.1:9876/.well-known/faramesh.json
curl -sf http://127.0.0.1:9876/v1/stats
```
