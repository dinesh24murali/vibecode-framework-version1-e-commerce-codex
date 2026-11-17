# Task: task-000-docker-toolchain

**Status:** done
**Created:** 2026-11-17
**ADR refs:** none

## Goal
Run every Go/Node/sqlc/goose/lint/e2e/terraform command inside Docker containers.

## Acceptance Criteria
- [x] docker-compose.tools.yml with go, node, e2e, tf services (profile tools)
- [x] .dockerignore for backend and frontend; .gitattributes forcing LF
- [x] Makefile targets wrap docker (go-test, go-lint, go-build, fe-build, fe-lint, migrate-*, sqlc-gen, codegen, verify, dev, seed)
- [x] verify scripts run in containers
- [x] docs/02_outputs/06_dev_setup.md and AGENTS.md updated
- [x] make go-build, go-test, go-lint, fe-build, migrate-status succeed

## Dependencies
- **Tasks:** none
- **Memory files:** none
- **Docs:** Ralph.md

## Files Expected to Change
docker-compose.tools.yml, Makefile, .gitattributes, .dockerignore files, verify/scripts/*, docs, AGENTS.md files
