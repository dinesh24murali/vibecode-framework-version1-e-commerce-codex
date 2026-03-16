# Intake Questionnaire

> Fill this out before using any prompt in `docs/01_prompts/`.
> Your answers replace the `[[PLACEHOLDER]]` tokens in each prompt.

---

## 1. Project Basics

**1.1 Project name:**
`[[PROJECT_NAME]]`

**1.2 One-sentence description:**
`[[PROJECT_DESCRIPTION]]`

**1.3 Team / company name:**
`[[TEAM_NAME]]`

**1.4 Target launch date (or "no fixed date"):**


---

## 2. Users & Problem

**2.1 Who are the primary users?**
(e.g. "Small business owners who need to track invoices")


**2.2 What problem does this solve?**


**2.3 What does success look like in 6 months?**
(Key metrics, user counts, revenue, etc.)


---

## 3. Core Features

**3.1 List the 3–5 must-have features for v1:**
1.
2.
3.
4.
5.

**3.2 List features explicitly OUT of scope for v1:**


**3.3 Are there any existing systems this must integrate with?**
(APIs, databases, legacy services)


---

## 4. Tech Stack

**4.1 Backend language:**
`[[BACKEND_LANG]]`  (e.g. Go, Python, Node.js, Ruby)

**4.2 Backend framework:**
`[[BACKEND_FRAMEWORK]]`  (e.g. Gin, FastAPI, Express, Rails)

**4.3 Frontend framework:**
`[[FRONTEND_FRAMEWORK]]`  (e.g. React, Vue, SvelteKit, Next.js)

**4.4 Database:**
`[[DATABASE]]`  (e.g. PostgreSQL, MySQL, SQLite, MongoDB)

**4.5 Cache / queue (if any):**
(e.g. Redis, RabbitMQ, none)

**4.6 Any preferred UI component library?**
(e.g. shadcn/ui, Tailwind, Material UI, none)

---

## 5. Authentication & Authorization

**5.1 Auth strategy:**
`[[AUTH_STRATEGY]]`  (e.g. JWT, OAuth2 + JWT, Session cookies, Passkeys)

**5.2 OAuth providers (if any):**
(e.g. Google, GitHub, none)

**5.3 Authorization model:**
(e.g. RBAC with roles: admin/user, simple owner-based, none)

---

## 6. Deployment & Infrastructure

**6.1 Deployment target:**
`[[DEPLOYMENT_TARGET]]`  (e.g. AWS ECS, GCP Cloud Run, Fly.io, Vercel + Railway)

**6.2 Containerized?**
(Yes — Docker / Docker Compose / Kubernetes | No)

**6.3 CI/CD preference:**
(e.g. GitHub Actions, GitLab CI, CircleCI, none)

**6.4 Environments needed:**
(e.g. dev, staging, prod)

---

## 7. Scale & Non-Functional Requirements

**7.1 Expected users at launch:**


**7.2 Expected peak requests/second:**


**7.3 Data sensitivity:**
(e.g. PII, financial data, public data only)

**7.4 Compliance requirements:**
(e.g. GDPR, HIPAA, SOC2, none)

**7.5 Uptime target:**
(e.g. 99.9%, best-effort)

---

## 8. Testing Expectations

**8.1 Minimum unit test coverage target:**
(e.g. 80%, none for now)

**8.2 E2e test tool preference:**
(e.g. Playwright, Cypress, none)

**8.3 Contract testing:**
(Yes — generate from OpenAPI spec | No)

---

## 9. Open Questions

List any unknowns or decisions not yet made:

1.
2.
3.
