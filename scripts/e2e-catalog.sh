#!/usr/bin/env bash
# Smoke test: catalog.json and a provider binary are reachable on GitHub raw.
set -euo pipefail
REF="${FARAMESH_REGISTRY_GITHUB_REF:-main}"
BASE="https://raw.githubusercontent.com/faramesh/faramesh-registry/${REF}"
curl -fsSL "${BASE}/catalog/catalog.json" | jq -e '.providers | length > 0' >/dev/null
curl -fsSL "${BASE}/catalog/artifacts/policies/faramesh/stripe/1.0.0/policy.fpl" | grep -q permit
curl -fsSL "${BASE}/catalog/artifacts/providers/faramesh/vault/1.0.0/manifest.json" | jq -e '.downloads.linux_amd64' >/dev/null
echo "e2e-catalog: GitHub raw artifacts OK (ref=${REF})"
