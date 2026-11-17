#!/usr/bin/env bash
# Installs Playwright test deps into /opt/tools (container-only). Version must match the e2e image.
set -euo pipefail
mkdir -p /opt/tools
cd /opt/tools
if [ ! -x node_modules/.bin/ts-node ]; then
  [ -f package.json ] || echo '{"name":"tools","private":true}' > package.json
  npm i --no-audit --no-fund playwright@1.50.0 @playwright/test@1.50.0 ts-node typescript@5.7.3 @types/node@22
fi
