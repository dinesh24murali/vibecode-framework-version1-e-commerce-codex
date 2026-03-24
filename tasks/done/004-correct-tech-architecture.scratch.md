# Scratch: 004-correct-tech-architecture

> AI: fill this out **before writing any code**. The developer will review this plan before implementation begins.
> Save as `tasks/active/004-correct-tech-architecture.scratch.md`.

---

## Plan

1. Update the technical architecture document's system topology to move static frontend hosting to S3 plus CloudFront and reduce environments to `dev`.
2. Correct the ERD and core data model sections for user roles, order address snapshots, product discounts and images, coupon tables, and the single-table cart model.
3. Add an extensible product attribute model for `level`, `theme`, `size`, and future filter dimensions without adding new columns to `products`.
4. Revise API pagination, caching, container/runtime notes, and the baseline summary to match the requested architecture.
5. Update persistent project memory and changelog entries, then archive the task files in `tasks/done/`.

## Files to Create or Modify

| File | Action | Notes |
|------|--------|-------|
| `tasks/active/004-correct-tech-architecture.md` | create | Task record for this documentation correction |
| `tasks/active/004-correct-tech-architecture.scratch.md` | create | Required scratch plan before implementation |
| `docs/02_outputs/03_tech_architecture.md` | modify | Apply the requested architecture and schema corrections |
| `memory/product.md` | modify | Persist product-domain schema decisions and extensible attribute model |
| `memory/project.md` | modify | Persist environment and frontend hosting decisions |
| `memory/_index.md` | modify | Refresh memory summaries after significant updates |
| `CHANGELOG.md` | modify | Add the required task completion entry |

## Uncertainties

- The original architecture assumed Nginx served the frontend; this update moves storefront delivery to S3 plus CloudFront and keeps API ingress details intentionally light.
- The `PRODUCT_IMAGE` and `CART` snippets provided mark multiple columns as `PK`; this update normalizes them to a single `id` primary key plus the intended foreign-key relationships.

## What Will Be Skipped

- No ADR will be created in this task because these are documentation corrections within the existing architecture output.
- No API spec or implementation code changes will be made here because the request is limited to the technical architecture and supporting memory/changelog records.

## Risks

- If later implementation expects the previous `CARTS` plus `CART_ITEMS` split, this documentation update intentionally changes that assumption.
- The extensible product attribute model introduces additional joins; that is the tradeoff for avoiding repeated schema changes when new dimensions are added.

## Ready to implement?

Yes. The plan and file list are defined and implementation can proceed.
