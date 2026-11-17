#!/usr/bin/env bash
# Runs check-dom.ts inside the e2e container (Playwright image). Tools are installed to /opt/tools.
set -euo pipefail
bash "$(dirname "$0")/install-e2e-deps.sh"
NODE_PATH=/opt/tools/node_modules /opt/tools/node_modules/.bin/ts-node --transpile-only --compiler-options '{"module":"commonjs","esModuleInterop":true,"lib":["es2022","dom"],"skipLibCheck":true,"target":"es2022"}' verify/scripts/check-dom.ts
