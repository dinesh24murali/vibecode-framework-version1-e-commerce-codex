# Memory Index

> Update this file whenever a memory file is created or significantly changed.
> AI tools: check this index before starting any task.

| File | Domain | Summary |
|------|--------|---------|
| `memory/project.md` | project | Stack, single `dev` environment, and S3 plus CloudFront frontend hosting notes |; frontend API client codegen (orval)
| `memory/product.md` | product | Storefront scope, product schema rules, cart and order-address behavior, coupon tracking, and extensible product attributes |
| `memory/infra.md` | infra | Docker-only toolchain, pgxpool limits, Redis AOF, migration workflow, sqlc codegen, local dev ports |; background worker (Redis Streams)
| `memory/auth.md` | auth | JWT (Ed25519), token lifetimes, Argon2id params, cookie settings, refresh token rotation, OTP delivery |; frontend Zustand store (tokens in memory, user persisted)

---

## How to Add an Entry

```
| `memory/<domain>.md` | <domain> | <one-line summary of what's in the file> |
```
