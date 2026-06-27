# Contributing to the Faramesh registry

Artifacts in this repository are signed policy packs, framework profiles, and provider binaries consumed by `faramesh check` and `faramesh apply`.

## Quick path (policy pack, ~30 minutes)

Full walkthrough with copy-paste commands: [Contribute to the registry](https://docs.faramesh.dev/registry/contributing/) on docs.faramesh.dev.

Summary:

1. Fork and clone this repo.
2. Copy `catalog/artifacts/policies/faramesh/stripe/1.0.0/` to `catalog/artifacts/policies/faramesh/<your-pack>/1.0.0/`.
3. Edit `policy.fpl` with your rules.
4. Add an entry to `catalog/catalog.json` with `"trust_tier": "community"`.
5. Run `./scripts/validate-catalog.sh`.
6. Open a pull request.

Faramesh uses the Developer Certificate of Origin (DCO). Every commit must
include a sign-off line:

```text
Signed-off-by: Your Name <you@example.com>
```

Use `git commit -s` to add the sign-off automatically.

## Layout

```
catalog/
  catalog.json
  trust/keys.json
  artifacts/
    policies/faramesh/<name>/<version>/policy.fpl
    frameworks/<name>/<version>/profile.fpl
    providers/faramesh/<name>/<version>/manifest.json
```

## New provider binary

1. Implement under `providers/<name>/` (see `providers/spiffe`).
2. Add to `Makefile` `PROVIDERS`.
3. `make providers` and `./scripts/refresh-provider-hashes.sh`.
4. Register in `catalog/catalog.json`.

Maintainers with `REGISTRY_SIGNING_KEY_B64` run `make sign-catalog`. See [SIGNING_KEY_LOCATION.md](./SIGNING_KEY_LOCATION.md) (maintainers only).

## Community tier

Set `"trust_tier": "community"` and document required config in a `README.md` beside the artifact. Consumers pin your publisher key in `governance.fms`.

## AI-Assisted Contributions

AI-assisted contributions are allowed when the contributor understands the
change, tests it, and takes responsibility for it. Fully automated, low-signal,
or undisclosed bulk AI-generated pull requests may be closed. Significant AI
assistance should be disclosed in the pull request body.

See [GITHUB_DISTRIBUTION.md](./GITHUB_DISTRIBUTION.md) for how stacks resolve artifacts from GitHub.
