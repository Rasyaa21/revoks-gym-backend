#!/usr/bin/env bash
set -euo pipefail

# Backward-compatible wrapper.
# Prefer using: ./scripts/run.sh <command>
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
exec "$SCRIPT_DIR/scripts/run.sh" "$@"
