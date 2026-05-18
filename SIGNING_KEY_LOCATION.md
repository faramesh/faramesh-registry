# Registry signing key

The Faramesh public registry signs every catalog manifest and provider binary with an Ed25519 keypair. Verifiers (the Faramesh CLI, the daemon, anyone running `validate-catalog.sh`) check signatures against the public key checked in at `catalog/trust/keys.json`.

This document is for **registry maintainers** — people who run `make sign-catalog` and rotate the signing key. Contributors and end users do not need a key.

## Where the key material lives

| Material | Path | Visibility |
|----------|------|------------|
| Private key (base64) | `.secrets/REGISTRY_SIGNING_KEY_B64` | Local-only; gitignored. Never commit. |
| Generation log | `.secrets/signing-key-generation.txt` | Local-only; metadata about the keypair. |
| Public key (current) | `catalog/trust/keys.json` → `faramesh-ed25519-2026` | Public; safe to commit and distribute. |

Public key fingerprint of the active keypair:

```
nxNjaQnS3L+zzKrRq48XfYBDWlFXkNJkxUiTD8j0sFs=
```

If `.secrets/REGISTRY_SIGNING_KEY_B64` is missing on a maintainer's machine, retrieve it from the team's secret store. It is never recreated except as part of an explicit key rotation.

## Setting the GitHub Actions secret

CI re-signs the catalog on every release. The base64-encoded private key is stored as the repo secret `REGISTRY_SIGNING_KEY_B64`.

To set or rotate the secret, from the **root of a `faramesh-registry` clone** (the directory containing `.secrets/`):

```bash
gh secret set REGISTRY_SIGNING_KEY_B64 \
  --repo faramesh/faramesh-registry \
  --body "$(cat .secrets/REGISTRY_SIGNING_KEY_B64)"
```

If you are running this from a parent monorepo folder, `cd faramesh-registry` first so the relative `.secrets/...` path resolves.

## Re-signing the catalog locally

```bash
cd faramesh-registry
export REGISTRY_SIGNING_KEY_B64="$(cat .secrets/REGISTRY_SIGNING_KEY_B64)"
make sign-catalog          # signs catalog manifests
make sign-all              # also rebuilds and signs provider binaries
./scripts/validate-catalog.sh
```

`validate-catalog.sh` confirms every signed artifact verifies against `catalog/trust/keys.json`.

## Rotating the key

1. Generate a new keypair (`make gen-signing-key`) and store the private half in `.secrets/REGISTRY_SIGNING_KEY_B64`.
2. Add the new public key to `catalog/trust/keys.json` alongside the old one.
3. Run `make sign-all` to produce signatures from the new key.
4. Verify with `./scripts/validate-catalog.sh`.
5. Open a PR with the updated `keys.json` and signatures.
6. After the PR merges and downstream consumers have picked up the new key, remove the old key entry from `keys.json` in a follow-up PR.
7. Update the GitHub Actions secret per the section above.
8. Update the public key fingerprint published in the docs (`faramesh-docs/content/docs/fpl.md`, security and concepts pages).

## Safety

- `.secrets/` is gitignored. Do not commit anything inside it.
- Never paste the private key into chat, issues, PRs, or CI logs.
- Only the public key in `catalog/trust/keys.json` is intended for distribution.
