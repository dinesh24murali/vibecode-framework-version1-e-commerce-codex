# Memory: Project

> Seed this file during project setup (task 000-project-setup).
> Update when the stack, team, or key decisions change.

---

## Project Identity

- **Name:** `e-commerce-site`
- **Description:** `This is a simple e-commerce-site for selling books`
- **Team:** `Feather Tech`
- **Repo:** `https://github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex`
- **Started:** `16-3-2026`

## Tech Stack

| Layer | Technology | Version |
|-------|-----------|---------|
| Backend lang | `Go 1.26` | |
| Backend framework | `Gin Web Framework` | |
| Frontend | `NextJS 16 Static Site Generation` | |
| Database | `Postgres` | |
| Auth | `JWT` | |
| Deployment | `AWS` | |
| CI/CD | `[[CI_CD_TOOL]]` | |

## Key Architectural Decisions

| Decision | ADR | Summary |
|----------|-----|---------|
| Use vibecode framework | ADR-0001 | File-based, AI-readable project structure |
| *(add more as ADRs are created)* | | |

## Environments

| Env | URL | Notes |
|-----|-----|-------|
| Dev | `http://localhost:3000` | Primary working environment via local Docker Compose |

## External Services

| Service | Purpose | Docs / Credentials |
|---------|---------|--------------------|
| Amazon S3 | Store the static frontend export | |
| Amazon CloudFront | Serve and cache the static frontend | |
| *(add as integrations are built)* | | |

## Key Conventions (Quick Reference)

- API contract: `docs/02_outputs/04_api_spec.yaml` is the source of truth
- Frontend client is **generated** from the API spec — never hand-written
- Static frontend deployments are uploaded manually to S3 and served through CloudFront in the current phase
- All dates stored as UTC ISO 8601
- *(add project-specific conventions here)*

## Open Questions

- *(add unresolved questions here; remove when answered)*
