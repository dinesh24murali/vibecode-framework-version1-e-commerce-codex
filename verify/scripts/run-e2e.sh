#!/usr/bin/env bash
# run-e2e.sh — Run end-to-end tests
#
# Usage: bash verify/scripts/run-e2e.sh [--headed] [--grep <pattern>]
#
# Environment variables:
#   APP_URL      Dev server URL (default: http://localhost:3000)
#   CI           Set to "true" in CI environments (enables retries, disables headed)

set -euo pipefail

APP_URL="${APP_URL:-http://localhost:3000}"
HEADED="${HEADED:-false}"
EXTRA_ARGS=()

# Parse args
while [[ $# -gt 0 ]]; do
  case "$1" in
    --headed)
      HEADED=true
      shift
      ;;
    --grep)
      EXTRA_ARGS+=("--grep" "$2")
      shift 2
      ;;
    *)
      EXTRA_ARGS+=("$1")
      shift
      ;;
  esac
done

echo "Running e2e tests against ${APP_URL}..."
echo ""

# Check that dev server is reachable
if ! curl -sf "${APP_URL}" > /dev/null 2>&1; then
  echo "ERROR: Dev server not reachable at ${APP_URL}"
  echo "Start it first, then re-run this script."
  exit 1
fi

# Build Playwright args
PLAYWRIGHT_ARGS=(
  "test"
  "--config=tests/e2e"
)

if [[ "${HEADED}" == "true" ]]; then
  PLAYWRIGHT_ARGS+=("--headed")
fi

if [[ "${CI:-false}" == "true" ]]; then
  PLAYWRIGHT_ARGS+=("--retries=2")
fi

if [[ ${#EXTRA_ARGS[@]} -gt 0 ]]; then
  PLAYWRIGHT_ARGS+=("${EXTRA_ARGS[@]}")
fi

# Run tests
APP_URL="${APP_URL}" npx playwright "${PLAYWRIGHT_ARGS[@]}"

EXIT_CODE=$?

if [[ $EXIT_CODE -eq 0 ]]; then
  echo ""
  echo "✓ All e2e tests passed"
else
  echo ""
  echo "✗ E2e tests failed (exit code ${EXIT_CODE})"
  echo "Check the Playwright HTML report: npx playwright show-report"
fi

exit $EXIT_CODE
