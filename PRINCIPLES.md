# Registry language principles

## FPL is primary

- Policy packs and framework profiles **must** ship canonical bytes in `policy.fpl` or `profile.fpl`.
- The CLI, compiler, and registry API treat FPL as the source of truth for install and checksums.
- Import strings in docs and UI always use FPL syntax:

```fpl
import "registry.faramesh.dev/policies/faramesh/stripe@1.0.0"
```

## YAML and JSON are supported naturally

- Publishers **may** attach `policy.yaml`, `policy.json`, `profile.yaml`, or `profile.json` sidecars.
- The API returns them when present; the web UI shows FPL first with YAML/JSON tabs.
- `GET ...?format=yaml|json` returns raw sidecar bytes (404 if not published).
- Governance stacks in user projects may be authored as FPL, YAML, or JSON — that is a **runtime** concern in `faramesh-core`. The **registry** still stores and distributes FPL-first artifacts.

## Providers

- Providers are not FPL; they are signed binaries with `manifest.json` and FPL **consumption** snippets in README.
