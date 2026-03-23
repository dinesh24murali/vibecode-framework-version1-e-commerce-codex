# Task — Generate Technical Architecture

- Task ID: `003`
- Status: `done`
- Type: `documentation`

## Goal
Create the production-ready technical architecture document for the initial release of `e-commerce-site`.

## Inputs
- Project name: `e-commerce-site`
- Backend: `Go 1.26 / Gin`
- Frontend: `NextJS 16 / SSG`
- Database: `Postgres`
- Deployment target: `AWS`
- Auth strategy: `JWT`
- Compliance: `DPDPA (India)`
- Containerized: `Docker Compose`
- CI/CD: `None`
- Initial deployment topology: `Single AWS EC2 instance running application and supporting services via Docker Compose`

## Deliverables
- Architecture document written to `docs/02_outputs/03_tech_architecture.md`
- Mermaid system, ERD, and auth sequence diagrams
- Opinionated production choices for infrastructure, observability, and security

## Notes
- Scratch planning captured in `tasks/done/003-generate-tech-architecture.scratch.md`
- Architecture direction updated to start with a simple single-EC2 deployment model for frontend, API, worker, PostgreSQL, and Redis via Docker Compose
- ADR candidates included in the architecture output for follow-up creation in `adr/`
