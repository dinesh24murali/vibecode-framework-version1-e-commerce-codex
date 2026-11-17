# Task: TASK-020 — Contact Us API + page

**Status:** done
**Created:** 2026-11-17
**ADR refs:** none

## Goal
Public contact form endpoint that stores enquiries and emails the business, plus the Contact Us page.

## Acceptance Criteria
- [x] POST /api/v1/contact stores the query and returns 200 + MessageResponse (the spec wins over the plan's 201)
- [x] Notification email sent via SMTP when CONTACT_NOTIFY_EMAIL is set (failure logged, enquiry still stored)
- [x] Google Maps iframe embedded for the configured address (NEXT_PUBLIC_BUSINESS_ADDRESS)
- [x] Form shows success on 200 and an error message on API failure
- [x] Client-side validation (required, email format, phone format)
- [x] Contract test backend/tests/contract/contact_test.go (200 shape, 400, 422, phone optional, stored + notified)
- [x] E2e test covers the happy path, validation and API failure
- [x] go build/test/lint, static build, e2e and DOM checks pass (docker)

## Dependencies
Tasks: 003, 005. Memory: memory/product.md, memory/infra.md.
