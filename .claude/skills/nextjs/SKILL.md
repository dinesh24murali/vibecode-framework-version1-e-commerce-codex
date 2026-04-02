---
name: nextjs
description: Project-specific Next.js patterns for the e-commerce storefront (SSG static export, App Router, Zustand, shadcn/ui, orval)
paths: "frontend/**/*.{ts,tsx,js,jsx}"
---

## Project Stack
- Next.js 16, TypeScript 5.7, App Router
- **Output: static export** (`output: 'export'`, `distDir: 'out'`) — deployed to S3 + CloudFront
- Zustand v5 for client state, shadcn/ui + Radix UI + Tailwind CSS for UI
- orval generates TypeScript API client from `docs/02_outputs/04_api_spec.yaml` → run `make codegen`

## Hard Rules
- **No SSR/ISR/Server Actions** — this is a fully static export. `getServerSideProps`, `revalidate`, and server-only APIs do not work.
- **Never manually write API clients** — regenerate from OpenAPI spec: `make codegen`
- **No `next/image` optimization** — images use `unoptimized: true` (CDN handles this)
- All imports use `@/` alias (maps to `frontend/`) — no relative `../../` imports

## App Router Structure
```
frontend/app/
├── (account)/    # Route group: account pages (no URL segment)
├── (admin)/      # Route group: admin pages
├── (auth)/       # Route group: auth/login pages
├── (shop)/       # Route group: storefront pages
├── layout.tsx    # Root layout — html/body/metadata only
└── page.tsx      # Home page
```

Route groups `(name)` organize pages without affecting URLs. Every new page needs:
1. The page file in the correct group
2. An e2e test in `tests/e2e/`
3. `make verify` must pass

## Zustand Store (ADR-0003 — SSG-safe hydration)
```typescript
// lib/store/index.ts — skipHydration: true prevents localStorage read during SSG build
const useStore = create(
  persist(storeSlice, {
    name: 'store',
    storage: createJSONStorage(() =>
      typeof window !== 'undefined' ? localStorage : memoryStorage
    ),
    skipHydration: true,   // REQUIRED for static export
  })
)

// In a client component root (e.g., layout or provider):
useEffect(() => { useStore.persist.rehydrate() }, [])
```
Never remove `skipHydration: true` — it will break the static build.

## Components
- UI primitives: `components/ui/` (shadcn/ui — generated, do not edit manually)
- Add new shadcn components via: `npx shadcn@latest add <component>`
- New feature components go in `components/<feature>/` outside `ui/`
- Icons: lucide-react

## Styling
- Tailwind CSS only — no inline styles, no CSS modules unless unavoidable
- Follow shadcn/ui's `cn()` utility from `lib/utils.ts` for conditional class merging

## Client vs Server Components
Since this is a static export, all interactive components need `'use client'`. Keep client boundaries as small as possible — wrap only the interactive leaf, not the whole page.

```typescript
// Good: small client boundary
'use client'
export function AddToCartButton({ productId }: { productId: string }) { ... }

// Bad: entire page as client component
'use client'
export default function ProductPage() { ... }
```

## Data Fetching Pattern (static export)
- Fetch at build time in page components using `generateStaticParams` for dynamic routes
- API calls at runtime go through the generated orval client (`lib/api/`)
- No `fetch` with `cache: 'no-store'` — incompatible with static export

## TypeScript
- Strict mode enabled — no `any`, no `@ts-ignore` without comment explaining why
- API types come from orval-generated client — do not manually define API response types
- Path alias: `@/*` → project root

## Verification
After any frontend change: `make verify` (runs Playwright e2e + DOM checks)
Dev server: `make dev` (runs on `http://localhost:3000`)
Build check: `cd frontend && npm run build` (must produce `out/` with no errors)
