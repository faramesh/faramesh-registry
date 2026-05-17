# Registry signing key — canonical location

**You do not need to hunt for a key.** The active Ed25519 keypair for this workspace is:

| Material | Path |
|----------|------|
| **Private key (base64)** | `faramesh-registry/.secrets/REGISTRY_SIGNING_KEY_B64` |
| **Generation log** | `faramesh-registry/.secrets/signing-key-generation.txt` |
| **Public key (in repo)** | `faramesh-registry/catalog/trust/keys.json` → `faramesh-ed25519-2026` |

Public key (verify only): `nxNjaQnS3L+zzKrRq48XfYBDWlFXkNJkxUiTD8j0sFs=`

## CI / GitHub Actions

Run from the **root of this repository** (where `.secrets/` lives):

```bash
gh secret set REGISTRY_SIGNING_KEY_B64 \
  --repo faramesh/faramesh-registry \
  --body "$(cat .secrets/REGISTRY_SIGNING_KEY_B64)"
```

If you are in a parent monorepo folder, do **not** prefix with `faramesh-registry/` — use `cat .secrets/...` only after `cd faramesh-registry`.

## Re-sign catalog locally

```bash
cd faramesh-registry   # repo root
export REGISTRY_SIGNING_KEY_B64="$(cat .secrets/REGISTRY_SIGNING_KEY_B64)"
make sign-catalog   # or make sign-all to rebuild provider bins too
./scripts/validate-catalog.sh
```

**Never commit** `.secrets/` (gitignored). The public key in `keys.json` is safe to commit.
