# ADR-0005 — OTP Delivery Mechanism

**Date:** 2026-04-03
**Status:** Accepted
**Task refs:** TASK-006

---

## Context

Both the registration flow and the password-reset flow require delivering a 6-digit OTP to the user's email address. A mechanism is needed to send (or simulate sending) this OTP. The architecture supports two distinct phases:

1. **Local / dev** — no email infrastructure is in place; a simple stdout log is sufficient.
2. **Production** — the OTP must be delivered reliably via email. Candidates are AWS SES (the target deployment is AWS) and self-hosted SMTP (MailHog in local dev, per the docker-compose setup).

## Decision

**Phase 1 (this task):** log the OTP to stdout using `slog.Info`. This is explicit, traceable in container logs, and requires no additional infrastructure.

**Phase 2 (future task):** replace the log call with an AWS SES send via the official Go SDK. MailHog (already in docker-compose on port 1025/8025) serves as the local SMTP relay for integration testing.

The OTP generation and Redis storage are unchanged between phases. The only change in Phase 2 is the delivery call inside `AuthService.RegisterInit` and `AuthService.PasswordResetInit`.

## Consequences

- No email is ever actually sent during Phase 1 development. Developers must observe the OTP in server logs.
- The code is structured so the delivery call is a single-line swap (stdout → SES), minimising the Phase 2 migration effort.
- MailHog is already provisioned in docker-compose for Phase 2 integration tests.

## Alternatives Considered

| Option | Reason rejected |
|--------|----------------|
| SMTP via MailHog from day one | Adds infra coupling before auth endpoints are proven. Deferred intentionally. |
| SendGrid | Project is AWS-first; SES is lower cost and avoids a third-party vendor. |
| Store OTP in Postgres instead of Redis | Redis with TTL is the simpler, more correct primitive for expiring short-lived codes. |
