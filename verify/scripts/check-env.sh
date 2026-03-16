#!/usr/bin/env bash
# check-env.sh — Validate that .env has all required variables from .env.example
#
# Usage: bash verify/scripts/check-env.sh [--env-file <path>] [--example-file <path>]
#
# Exit codes:
#   0 — all required vars are set
#   1 — one or more required vars are missing or empty

set -euo pipefail

ENV_FILE="${ENV_FILE:-.env}"
EXAMPLE_FILE="${EXAMPLE_FILE:-.env.example}"

# Parse args
while [[ $# -gt 0 ]]; do
  case "$1" in
    --env-file)     ENV_FILE="$2";     shift 2 ;;
    --example-file) EXAMPLE_FILE="$2"; shift 2 ;;
    *) echo "Unknown arg: $1"; exit 1 ;;
  esac
done

if [[ ! -f "$EXAMPLE_FILE" ]]; then
  echo "ERROR: Example file not found: $EXAMPLE_FILE"
  exit 1
fi

if [[ ! -f "$ENV_FILE" ]]; then
  echo "ERROR: .env file not found at: $ENV_FILE"
  echo ""
  echo "Create it with:"
  echo "  cp $EXAMPLE_FILE $ENV_FILE"
  echo "Then fill in all required values."
  exit 1
fi

# Extract required variable names from .env.example:
# - Skip comment lines (starting with #)
# - Skip blank lines
# - Skip lines where the value is already set (KEY=value) — only flag empty (KEY=)
required_vars=()
while IFS= read -r line; do
  # Skip comments and blank lines
  [[ "$line" =~ ^[[:space:]]*# ]] && continue
  [[ -z "${line// }" ]] && continue

  # Extract key (everything before first =)
  key="${line%%=*}"
  key="${key// /}"  # trim spaces
  [[ -z "$key" ]] && continue

  required_vars+=("$key")
done < "$EXAMPLE_FILE"

# Check each required var in .env
missing=()
for var in "${required_vars[@]}"; do
  # Source the var from .env (safe: only read, no execute)
  value=$(grep -E "^${var}=" "$ENV_FILE" 2>/dev/null | head -1 | cut -d'=' -f2-)
  if [[ -z "$value" ]]; then
    missing+=("$var")
  fi
done

if [[ ${#missing[@]} -eq 0 ]]; then
  echo "✓ All required environment variables are set (${#required_vars[@]} checked)"
  exit 0
else
  echo "✗ Missing or empty environment variables in $ENV_FILE:"
  echo ""
  for var in "${missing[@]}"; do
    echo "  - $var"
  done
  echo ""
  echo "Set these in $ENV_FILE and re-run:"
  echo "  bash verify/scripts/check-env.sh"
  exit 1
fi
