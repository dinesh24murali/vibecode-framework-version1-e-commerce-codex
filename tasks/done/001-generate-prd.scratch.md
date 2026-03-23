# Scratch: 001-generate-prd

## Plan

1. Review repository instructions for docs, tasks, and memory handling.
2. Draft a production-ready PRD that covers storefront, account, checkout, admin, and platform requirements.
3. Save the PRD to `docs/02_outputs/01_prd.md` using the repository naming convention.
4. Record durable product scope facts in project memory for future implementation tasks.
5. Update `CHANGELOG.md`, mark the task done, and archive the task files into `tasks/done/`.

## Files to Change

- `tasks/active/001-generate-prd.md`
- `tasks/active/001-generate-prd.scratch.md`
- `docs/02_outputs/01_prd.md`
- `memory/_index.md`
- `memory/product.md`
- `CHANGELOG.md`

## Uncertainties

- Whether the product should support guest checkout or require authentication before checkout.
- Whether Indian GST is the expected tax regime, inferred from GST Number and Pincode fields.
- Whether the admin site is part of the same frontend deployment or a separate NextJS route group.

## What Will Be Skipped

- UI mockups, wireframes, and visual design specifications beyond textual requirements.
- Technical implementation details that belong in architecture or API design documents rather than the PRD.
- Payment gateway requirements, because they are explicitly out of scope for v1.

## Risks

- Some navigation labels are inconsistent (`Book Store` vs `Book Stores`, `Printing service` vs `Printing Services`), so assumptions must be called out explicitly.
- The input does not define inventory sync, shipping workflow, or order taxation rules in detail, so the PRD must make pragmatic assumptions to avoid ambiguity.
