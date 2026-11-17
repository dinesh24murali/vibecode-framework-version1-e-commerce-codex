# Scratch: TASK-020

## Plan
1. Repo (sqlc CreateContactQuery exists) -> ContactService (validation, notifier interface, SMTP notifier) -> handler -> router + config.
2. Contract test with fake repo/notifier. 3. Frontend ContactForm, page, business/social constants, e2e.

## Files
backend: repositories/contact_repository.go, services/contact_service.go, handlers/contact_handler.go, config.go, router.go, tests/contract/contact_test.go; frontend: app/(shop)/contact/page.tsx, components/shop/ContactForm.tsx, lib/constants/business.ts; tests/e2e/contact.spec.ts; .env.example.

## Decisions
- Spec says 200 and requires `subject`; plan said 201 / no subject. Spec wins; form has a Subject field.
- Business address is not in the docs: env-configured with a placeholder default.
- Email failure does not fail the request. Header values are stripped of CR/LF.
- No rate limiting (not in spec); noted as a risk for spam.

## Skipped
Admin listing of enquiries; CAPTCHA.

## Risks
Public unauthenticated write endpoint can be spammed; add rate limiting later.
