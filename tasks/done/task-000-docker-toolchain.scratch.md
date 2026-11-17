# Scratch: task-000-docker-toolchain

## Plan
1. Fix postgres 18 volume mount (new volume postgres18_data at /var/lib/postgresql).
2. Create docker-compose.tools.yml (go, node, e2e, tf) with named cache volumes.
3. Add .gitattributes, .dockerignore files.
4. Rewrite Makefile targets to docker.
5. Update verify scripts and docs.
6. Prove with make go-build/test/lint, fe-build, migrate-status.

## Files
docker-compose.yml, docker-compose.tools.yml, Makefile, .gitattributes, backend/.dockerignore, frontend/.dockerignore, verify/scripts/run-e2e.sh, docs/02_outputs/06_dev_setup.md, AGENTS.md files

## Uncertainties
- Decision: go image is custom-built (golang:1.26 + golangci-lint + sqlc + goose) via docker/go.Dockerfile.

## Skipped
- Real e2e tests (none exist yet).

## Risks
- Image pulls/golangci-lint version compatibility with Go 1.26.
