# Scratch: TASK-003 — Frontend scaffold — NextJS 16 App Router

## Plan

1. Write `frontend/package.json` manually (pinned to `next@16.2.1`, React 19, TypeScript, Tailwind, Zustand v5, shadcn/ui deps) to avoid interactive `create-next-app` prompts
2. Run `npm install` inside `frontend/`
3. Write `frontend/next.config.ts` with `output: 'export'` and `distDir: 'out'`
4. Write `frontend/tsconfig.json`
5. Write Tailwind config (`tailwind.config.ts`, `postcss.config.mjs`)
6. Write root `app/layout.tsx`, `app/page.tsx`, `app/globals.css`
7. Create the four route-group shell layouts: `(shop)`, `(auth)`, `(account)`, `(admin)`
8. Initialize shadcn/ui via `npx shadcn@latest init --yes --defaults` (or write `components.json` + button component manually)
9. Create Zustand store skeleton in `frontend/lib/store/index.ts` using `skipHydration` — SSG-safe pattern
10. Create `frontend/lib/store/slices/.gitkeep`
11. Run `npm run build` inside `frontend/` — verify `out/` is produced with no errors
12. Create ADR `adr/0002-zustand-ssr-hydration-strategy.md`
13. Update `CHANGELOG.md`
14. Move task + scratch to `tasks/done/`

## Files to Change / Create

- `frontend/package.json` *(create)*
- `frontend/package-lock.json` *(create — after npm install)*
- `frontend/tsconfig.json` *(create)*
- `frontend/next.config.ts` *(create — output: 'export', distDir: 'out')*
- `frontend/tailwind.config.ts` *(create)*
- `frontend/postcss.config.mjs` *(create)*
- `frontend/app/globals.css` *(create)*
- `frontend/app/layout.tsx` *(create — root layout)*
- `frontend/app/page.tsx` *(create — minimal index)*
- `frontend/app/(shop)/layout.tsx` *(create)*
- `frontend/app/(auth)/layout.tsx` *(create)*
- `frontend/app/(account)/layout.tsx` *(create)*
- `frontend/app/(admin)/layout.tsx` *(create)*
- `frontend/lib/store/index.ts` *(create)*
- `frontend/lib/store/slices/.gitkeep` *(create)*
- `frontend/components/ui/button.tsx` *(create — shadcn sample component)*
- `frontend/components.json` *(create — shadcn config)*
- `adr/0002-zustand-ssr-hydration-strategy.md` *(create)*
- `CHANGELOG.md` *(update)*

## Uncertainties

- shadcn/ui non-interactive init: will use `--yes` / `--defaults` if supported, otherwise write `components.json` and the button component files directly
- `npx shadcn@latest` version compatibility with Next 16.2.1 — expected to work; shadcn tracks Next releases closely

## What Will Be Skipped

- No feature pages (shop listing, product detail, auth forms, account pages) — later tasks
- No API client codegen (TASK-007)
- No real Zustand slices beyond skeleton (TASK-008)
- No e2e test (no feature page yet; verify baseline is just a successful build)

## Risks

- `output: 'export'` requires all pages to be statically renderable — shell layouts with no data fetching are safe
- shadcn/ui writes to `components/ui/` by default; must ensure `components.json` `aliases` match the project layout
