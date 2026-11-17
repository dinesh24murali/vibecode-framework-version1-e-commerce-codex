#!/usr/bin/env bash
# Fails if a raw fetch( or axios( call exists in frontend/ outside lib/api/client.ts and generated code.
set -e
hits=$(grep -rnE "\b(fetch|axios)\(" frontend --include=*.ts --include=*.tsx \
  --exclude-dir=node_modules --exclude-dir=.next --exclude-dir=out --exclude-dir=generated \
  | grep -v "frontend/lib/api/client.ts" || true)
if [ -n "$hits" ]; then echo "Raw fetch/axios calls found:"; echo "$hits"; exit 1; fi
echo "No raw fetch/axios calls outside lib/api/client.ts"
