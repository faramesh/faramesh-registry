# Contributing registry artifacts

## Layout

```
catalog/
  catalog.json              # index (required)
  trust/keys.json           # official Ed25519 public keys
  artifacts/
    providers/faramesh/<name>/<version>/manifest.json + bin/
    policies/faramesh/<name>/<version>/policy.fpl
    frameworks/<name>/<version>/profile.fpl
```

## New policy pack or framework profile

1. Copy an existing version directory (e.g. `policies/faramesh/stripe/1.0.0/`).
2. Edit canonical `policy.fpl` or `profile.fpl` (FPL is primary).
3. Add an entry to `catalog/catalog.json`.
4. Run `./scripts/validate-catalog.sh`.
5. If you have `REGISTRY_SIGNING_KEY_B64`, run `go run ./cmd/sign-catalog -catalog catalog`.

## New provider (binary)

1. Implement `ProviderService` under `providers/<name>/` (see `providers/vault`).
2. Add to `Makefile` `PROVIDERS` list.
3. `make providers && make sign-catalog -providers` (with signing key).
4. Register in `catalog/catalog.json` with capabilities.

## Community tier

Set `"trust_tier": "community"` and document required config in `README.md` beside the artifact. Users must pin your publisher key in their stack `trust` block.

See [GITHUB_DISTRIBUTION.md](./GITHUB_DISTRIBUTION.md) for how consumers resolve artifacts without registry.faramesh.dev.
