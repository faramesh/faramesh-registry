#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
CATALOG="$ROOT/catalog/catalog.json"
TRUST="$ROOT/catalog/trust/keys.json"

die() { echo "validate-catalog: $*" >&2; exit 1; }

command -v jq >/dev/null 2>&1 || die "jq is required"
[ -f "$CATALOG" ] || die "missing $CATALOG"
[ -f "$TRUST" ] || die "missing $TRUST"

check_versions() {
  local section="$1"
  local primary="$2"
  jq -r ".$section[]? | \"\(.name)|\(.latest_version)\"" "$CATALOG" | while IFS='|' read -r name latest; do
    [ -n "$name" ] || continue
    jq -r ".$section[] | select(.name==\"$name\") | .versions | keys[]" "$CATALOG" | while read -r ver; do
      path=$(jq -r ".$section[] | select(.name==\"$name\") | .versions[\"$ver\"]" "$CATALOG")
      rel="$ROOT/catalog/$path"
      if [ "$section" = "providers" ]; then
        [ -f "$rel" ] || die "provider manifest missing: $rel"
      else
        dir="$(dirname "$rel")"
        fpl="$dir/$primary"
        [ -f "$fpl" ] || die "canonical FPL missing: $fpl ($section/$name@$ver)"
        if [ -f "$fpl.sig" ] && command -v go >/dev/null 2>&1; then
          :
        fi
      fi
    done
    [ "$latest" != "" ] || die "$section/$name missing latest_version"
  done
}

check_versions policies policy.fpl
check_versions frameworks profile.fpl
check_versions providers manifest.json
check_versions packs policy.fpl

# Ensure demo pack has FPL primary (not YAML-only)
demo_fpl="$ROOT/catalog/artifacts/policies/faramesh/demo/0.1.0/policy.fpl"
[ -s "$demo_fpl" ] || die "demo policy.fpl must be non-empty (FPL is primary)"

echo "catalog OK: $CATALOG (FPL-primary artifacts validated)"
