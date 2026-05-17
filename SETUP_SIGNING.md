# Registry signing setup

## Public key vs private key

| Material | Where it lives | Safe to publish? |
|----------|----------------|------------------|
| **Public key** (`public_key_b64` in `catalog/trust/keys.json`) | Repo + docs `trust { ... }` block | **Yes** — required for verification |
| **Private key** (`REGISTRY_SIGNING_KEY_B64`) | GitHub Actions secret only | **Never** commit or paste in chat |

The line in `GITHUB_DISTRIBUTION.md` is the **public** key. That is correct.

## One-time: add CI secret

You need the Ed25519 private key that matches `faramesh-ed25519-2026` in `catalog/trust/keys.json`.

**If you still have the key** from when the catalog was first signed:

```bash
gh secret set REGISTRY_SIGNING_KEY_B64 --repo faramesh/faramesh-registry --body "$REGISTRY_SIGNING_KEY_B64"
```

**If the private key was lost**, rotate:

```bash
cd faramesh-registry
go run ./cmd/gen-signing-key
# 1. Save REGISTRY_SIGNING_KEY_B64 output as the GitHub secret (command above)
# 2. Replace catalog/trust/keys.json with the printed public entry
# 3. Re-sign everything:
export REGISTRY_SIGNING_KEY_B64="<paste private b64>"
make sign-all
git add catalog && git commit -m "Re-sign catalog after key rotation"
```

## What CI does when the secret is set

On push to `catalog/**`, `.github/workflows/catalog.yml`:

1. `./scripts/validate-catalog.sh`
2. `go run ./cmd/sign-catalog -catalog catalog` (FPL + profiles)
3. `go run ./cmd/sign-catalog -catalog catalog -providers` (binaries + manifests)

Without the secret, validation still runs; signing is skipped.

## Local verify

```bash
./scripts/validate-catalog.sh
export REGISTRY_SIGNING_KEY_B64="..."
make sign-all
```
