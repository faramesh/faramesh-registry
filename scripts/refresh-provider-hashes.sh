#!/usr/bin/env bash
# Updates sha256_hex and size in provider manifest.json from bin/ (no Ed25519 signing).
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
command -v jq >/dev/null || { echo "jq required"; exit 1; }

for man in "$ROOT"/catalog/artifacts/providers/faramesh/*/1.0.0/manifest.json; do
  dir="$(dirname "$man")"
  tmp="$(mktemp)"
  jq -c '.downloads | keys[]' "$man" | tr -d '"' | while read -r plat; do
    bin_rel=$(jq -r ".downloads[\"$plat\"].url" "$man" | sed 's|^file://||')
    bin="$dir/$bin_rel"
    [ -f "$bin" ] || { echo "missing $bin" >&2; exit 1; }
    hash=$(shasum -a 256 "$bin" | awk '{print $1}')
    size=$(wc -c <"$bin" | tr -d ' ')
    jq --arg p "$plat" --arg h "$hash" --argjson s "$size" \
      '.downloads[$p].sha256_hex = $h | .downloads[$p].size = $s' "$man" >"$tmp" && mv "$tmp" "$man"
    echo "  $plat $hash"
  done
  echo "updated $(basename "$(dirname "$dir")") manifest"
done
