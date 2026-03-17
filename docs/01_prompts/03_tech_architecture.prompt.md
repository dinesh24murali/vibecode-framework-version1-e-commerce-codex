# Prompt: Generate Technical Architecture Document

> **How to use:** Complete Prompts 01 and 02 first. Fill in every `[[PLACEHOLDER]]`, then paste this entire file into your AI tool. Save the output to `docs/02_outputs/03_tech_architecture.md`.

---

## Instructions to AI

You are a principal engineer. Generate a **production-ready technical architecture document** for the project below. Make concrete, opinionated technology choices. Do not hedge with "you could use X or Y" — pick one and justify it.

---

## Project Input

**Project name:** e-commerce-site
**Backend lang / framework:** Go 1.26 / Gin Web Framework
**Frontend:** NextJS 16
**Database:** Postgres
**Deployment target:** AWS
**Auth strategy:** JWT
**Scale:** [[Fill from intake 7.1–7.2]]
**Compliance:** [[Fill from intake 7.4]]
**Containerized:** [[Fill from intake 6.2]]
**CI/CD:** [[Fill from intake 6.3]]

---

## Required Sections

1. **System Overview** — high-level diagram (Mermaid) showing all components and their relationships
2. **Component Breakdown** — for each component (frontend, backend, DB, cache, queue, CDN, etc.):
   - Technology choice with version
   - Responsibility
   - Interfaces (what it exposes, what it depends on)
3. **Data Architecture**
   - Entity-relationship diagram (Mermaid ERD)
   - Core data models with field names and types
   - Indexing strategy
   - Migration strategy
4. **API Design Principles** — REST vs GraphQL vs RPC decision, versioning strategy, pagination pattern
5. **Authentication & Authorization Flow** — sequence diagram (Mermaid) of the JWT flow
6. **Infrastructure & Deployment**
   - Environment topology (dev / staging / prod)
   - Container / serverless configuration
   - Secrets management approach
   - Deployment pipeline overview
7. **Observability** — logging, metrics, tracing, alerting choices
8. **Security Considerations** — threat model highlights, input validation, rate limiting, CORS, CSP
9. **Scalability & Performance** — caching strategy, DB connection pooling, CDN usage, bottleneck analysis
10. **ADR Candidates** — list 3–5 decisions that should become ADRs (to be created in `adr/`)

---

## Output Format

- Use Markdown
- Use Mermaid for all diagrams
- Be specific with versions (e.g. "PostgreSQL 16", "React 18.3")
- Target length: 2500–5000 words
