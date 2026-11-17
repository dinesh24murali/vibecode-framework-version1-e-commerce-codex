# Scratch: TASK-018

## Plan
1. lib/constants/social.ts: social URLs from NEXT_PUBLIC_FACEBOOK_URL / INSTAGRAM_URL / WHATSAPP_URL (defaults to platform home pages); documented in .env.example. Reused by the contact page (TASK-020).
2. components/layout/Footer.tsx (server component; lucide icons, WhatsApp uses MessageCircle since lucide has no brand icon).
3. Mount in app/(shop)/layout.tsx. e2e tests/e2e/footer.spec.ts.

## Files
frontend/lib/constants/social.ts, frontend/components/layout/Footer.tsx, frontend/app/(shop)/layout.tsx, .env.example, tests/e2e/footer.spec.ts, CHANGELOG, PROGRESS.

## Decisions
Real social handles are not specified in the docs; env-configured with generic defaults.

## Skipped
Target pages (TASK-019). Footer on (auth)/(account) layouts: spec says public pages.

## Risks
Links 404 until TASK-019 lands.
