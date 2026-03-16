# Memory: Project

> Seed this file during project setup (task 000-project-setup).
> Update when the stack, team, or key decisions change.

---

## Project Identity

- **Name:** `[[PROJECT_NAME]]`
- **Description:** `[[PROJECT_DESCRIPTION]]`
- **Team:** `[[TEAM_NAME]]`
- **Repo:** `[[REPO_URL]]`
- **Started:** `[[PROJECT_START_DATE]]`

## Tech Stack

| Layer | Technology | Version |
|-------|-----------|---------|
| Backend lang | `[[BACKEND_LANG]]` | |
| Backend framework | `[[BACKEND_FRAMEWORK]]` | |
| Frontend | `[[FRONTEND_FRAMEWORK]]` | |
| Database | `[[DATABASE]]` | |
| Auth | `[[AUTH_STRATEGY]]` | |
| Deployment | `[[DEPLOYMENT_TARGET]]` | |
| CI/CD | `[[CI_CD_TOOL]]` | |

## Key Architectural Decisions

| Decision | ADR | Summary |
|----------|-----|---------|
| Use vibecode framework | ADR-0001 | File-based, AI-readable project structure |
| *(add more as ADRs are created)* | | |

## Environments

| Env | URL | Notes |
|-----|-----|-------|
| Local | `http://localhost:[[APP_PORT]]` | |
| Staging | `[[STAGING_URL]]` | |
| Production | `[[PROD_URL]]` | |

## External Services

| Service | Purpose | Docs / Credentials |
|---------|---------|--------------------|
| *(add as integrations are built)* | | |

## Key Conventions (Quick Reference)

- API contract: `docs/02_outputs/04_api_spec.yaml` is the source of truth
- Frontend client is **generated** from the API spec — never hand-written
- All dates stored as UTC ISO 8601
- *(add project-specific conventions here)*

## Open Questions

- *(add unresolved questions here; remove when answered)*
