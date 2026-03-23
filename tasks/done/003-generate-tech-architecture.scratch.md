# Scratch Notes — Generate Technical Architecture

## Objective
Produce a production-ready technical architecture document for `e-commerce-site` using the project constraints:
- Backend: Go 1.26 + Gin
- Frontend: NextJS 16 with static generation
- Database: PostgreSQL
- Deployment: AWS
- Auth: JWT
- Compliance: DPDPA (India)
- Local environment: Docker Compose

## Plan
1. Define the target production topology and service boundaries.
2. Choose concrete AWS services, runtime versions, and operational tooling.
3. Design the core data model for catalog, cart, order, payment, and identity domains.
4. Specify API style, versioning, authn/authz, and pagination standards.
5. Document deployment, observability, security, and scalability patterns.
6. Capture follow-up ADR candidates.

## Open Questions
- Payment gateway is not specified; architecture assumes Razorpay as the primary India-focused PSP.
- Search requirements are not explicit; architecture keeps PostgreSQL full-text search for the initial scale target instead of introducing OpenSearch.
- Staging is included even though CI/CD is currently absent; release execution is manual first, then automated later.

## Skips
- No endpoint-by-endpoint OpenAPI contract here; that belongs in the API spec output.
- No Terraform module breakdown; this document stays at architecture level.
- No cost model; the current deliverable is technical architecture.

## Estimated Files
- `docs/02_outputs/03_tech_architecture.md`
- `tasks/done/003-generate-tech-architecture.md`
- `tasks/done/003-generate-tech-architecture.scratch.md`
