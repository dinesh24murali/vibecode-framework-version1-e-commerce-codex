# Scratch: TASK-019

## Plan
1. Shared `components/shop/StaticPage.tsx` (page shell, Section, PlaceholderImage).
2. Four server-component pages under app/(shop)/ with `metadata` titles.
3. e2e static-pages.spec.ts; build + verify.

## Files
frontend/components/shop/StaticPage.tsx, frontend/app/(shop)/{wall-of-fame,about,terms,privacy}/page.tsx, tests/e2e/static-pages.spec.ts, CHANGELOG, PROGRESS.

## Decisions
Placeholder images are styled divs (no binary assets or external hosts). Legal copy is boilerplate flagged for counsel review. Privacy link is not in the footer spec, so the footer is unchanged.

## Skipped
Real content; footer link to privacy.

## Risks
None significant.
