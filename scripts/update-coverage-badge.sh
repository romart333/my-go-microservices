#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
COVERAGE_FILE="${1:-$ROOT_DIR/coverage/coverage.out}"
BADGE_FILE="${2:-$ROOT_DIR/coverage/coverage.json}"

if [[ ! -f "$COVERAGE_FILE" ]]; then
  echo "Coverage file not found: $COVERAGE_FILE"
  echo "Run: task test:coverage"
  exit 1
fi

if [[ -z "${COVERAGE_GIST_ID:-}" ]]; then
  echo "COVERAGE_GIST_ID is required"
  echo "Export it: export COVERAGE_GIST_ID=<gist_id>"
  exit 1
fi

TOKEN="${GIST_TOKEN:-${GITHUB_TOKEN:-}}"
if [[ -z "${TOKEN}" ]]; then
  echo "GIST_TOKEN (or GITHUB_TOKEN) is required"
  echo "Export it: export GIST_TOKEN=<token_with_gist_scope>"
  exit 1
fi

COVERAGE_PERCENT="$(go tool cover -func="$COVERAGE_FILE" | awk '/^total:/{print $3}')"

mkdir -p "$(dirname "$BADGE_FILE")"
cat >"$BADGE_FILE" <<EOF
{
  "schemaVersion": 1,
  "label": "coverage",
  "message": "$COVERAGE_PERCENT",
  "color": "green"
}
EOF

CONTENT_ESCAPED="$(sed ':a;N;$!ba;s/\\/\\\\/g;s/"/\\"/g;s/\n/\\n/g' "$BADGE_FILE")"
PAYLOAD="{\"files\":{\"coverage.json\":{\"content\":\"${CONTENT_ESCAPED}\"}}}"

curl --fail --silent --show-error \
  -X PATCH \
  -H "Accept: application/vnd.github+json" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  "https://api.github.com/gists/${COVERAGE_GIST_ID}" \
  --data "$PAYLOAD" >/dev/null

echo "Badge updated: $COVERAGE_PERCENT"
echo "Gist: https://gist.github.com/${COVERAGE_GIST_ID}"
