# Task: TASK-003 — Frontend scaffold — NextJS 16 App Router

**Status:** active
**Created:** 2026-03-31
**ADR refs:** ADR needed — Zustand SSR/SSG hydration strategy

---

## Goal

Bootstrap NextJS 16 with App Router, TypeScript, Tailwind CSS, and shadcn/ui; configure static export mode; install Zustand v5 and create the slice-pattern store skeleton; create all four route-group shell directories.

## Background

All frontend feature tasks (TASK-007 onward) depend on this scaffold. The `output: 'export'` configuration is mandatory — the storefront is a static site deployed to S3/CloudFront, not a Node.js server. The Zustand store must be safe to initialise in an SSG context (no `localStorage` access during the build).

See `docs/02_outputs/03_tech_architecture.md` §2 (Frontend) and `docs/02_outputs/05_implementation_plan.md` TASK-003.

## Acceptance Criteria

- [ ] `next build` produces a static export in `out/` with no errors
- [ ] `next.config.ts` has `output: 'export'` configured
- [ ] All four route-group directories exist with shell `layout.tsx` files: `app/(shop)/`, `app/(auth)/`, `app/(account)/`, `app/(admin)/`
- [ ] Zustand v5 store initialises without errors in SSG context — no `localStorage` access during static generation
- [ ] Tailwind CSS and shadcn/ui are installed and resolvable (a sample shadcn component renders without error)
- [ ] `make verify` passes (baseline — no feature pages yet, just scaffold)
- [ ] ADR created documenting the Zustand SSR/SSG hydration strategy

## Dependencies

- **Tasks:** TASK-001 (monorepo must exist)
- **Memory files:** none (first frontend task)
- **Docs:** `docs/02_outputs/03_tech_architecture.md` §2 (Frontend), `frontend/AGENTS.md`, `docs/02_outputs/05_implementation_plan.md` TASK-003

## Files Expected to Change

- `frontend/package.json` *(create)*
- `frontend/tsconfig.json` *(create)*
- `frontend/next.config.ts` *(create)*
- `frontend/tailwind.config.ts` *(create)*
- `frontend/app/layout.tsx` *(create)*
- `frontend/app/page.tsx` *(create)*
- `frontend/app/(shop)/layout.tsx` *(create)*
- `frontend/app/(auth)/layout.tsx` *(create)*
- `frontend/app/(account)/layout.tsx` *(create)*
- `frontend/app/(admin)/layout.tsx` *(create)*
- `frontend/lib/store/index.ts` *(create)*
- `frontend/lib/store/slices/` *(create directory with `.gitkeep`)*
- `adr/NNNN-zustand-ssr-hydration-strategy.md` *(create)*

## Notes

- NextJS version: `16.0.x`; React: `19.x`; Node.js build runtime: `22.x` (per architecture doc)
- UI library: `shadcn/ui`
- Zustand version: v5
- `output: 'export'` means no Server Actions, no `getServerSideProps`, no API routes in Next.js — all data must come from the backend API
- The Zustand SSG safety pattern typically involves checking `typeof window !== 'undefined'` before accessing `localStorage`, or using the `skipHydration` option — document the chosen approach in the ADR
- Do not implement any feature pages in this task — only shell layouts

---

## Scratch File

AI: before implementing, create `tasks/active/task-003-frontend-scaffold-nextjs.scratch.md` with your plan.
See `tasks/AGENTS.md` for the required scratch file format.
