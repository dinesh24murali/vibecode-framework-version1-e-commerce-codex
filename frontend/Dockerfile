# ── Build stage ───────────────────────────────────────────────────────────────
FROM node:22-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json* ./
RUN npm ci

COPY . .

# NEXT_PUBLIC_* vars are inlined at build time
ARG NEXT_PUBLIC_API_URL=http://localhost:8080
ENV NEXT_PUBLIC_API_URL=$NEXT_PUBLIC_API_URL

# Build static export
RUN npm run build

# ── Runtime stage ─────────────────────────────────────────────────────────────
# Static export is served via a simple HTTP server
FROM node:22-alpine

RUN npm install -g serve

WORKDIR /app

COPY --from=builder /app/out ./out

EXPOSE 3000

# No -s: a Next static export has one HTML file per route; SPA fallback would serve the home page everywhere
CMD ["serve", "out", "-l", "3000"]
