# Progress

Order: 000, 008, 007, 040, 009 ... 038. Status: todo/done/blocked.

| Task | Status | Note |
|---|---|---|
| 000 | done | make is not installed on host; run the docker compose equivalents (see AGENTS.md). Fixed 4 lint issues, removed debug prints. |
| 008 | done | CartApi adapter injected; TASK-007 wires generated client. |
| 007 | done | orval, generated client committed. |
| 040 | done | |
| 009 | done | e2e mocks API; Facebook popup not e2e-tested. |
| 010 | done | Phone not in spec; reload requires re-login (no refresh flow yet). |
| 011 | done | Filter predicate factory deferred to 012. |
| 012 | done | 12 seeded categories only after TASK-030; predicate factory not needed (SQL params). |
| 013 | done | Not tested against real AWS; presigned GET not wired (spec uses CDN public_url). |
| 014 | done | Product link target pages and slide admin come later; slides need seed (030). |
| 015 | done | Card links to product detail await TASK-016. |
| 016 | done | Slug->id resolved via list API (spec has only GET by id); fallback slug sample-product when API down at build. |
| 017 | done | Uses category slug + q params (matching TASK-015). Destination pages come in 018-020. |
| 018 | done | Social URLs env-configured (NEXT_PUBLIC_*_URL); links to terms/about 404 until 019. |
| 019 | done | Placeholder content; legal copy is boilerplate needing counsel review. |
| 020 | done | Spec wins: 200 + subject field. Address and notify email env-configured; no rate limiting yet. |
| 021 | todo | |
| 022 | todo | |
| 023 | todo | |
| 039 | todo | |
| 024 | todo | |
| 025 | todo | |
| 026 | todo | |
| 027 | todo | |
| 028 | todo | |
| 029 | todo | |
| 041 | todo | |
| 030 | todo | |
| 031 | todo | |
| 032 | todo | |
| 033 | todo | |
| 034 | todo | |
| 035 | todo | |
| 036 | todo | |
| 037 | todo | |
| 038 | todo | |
