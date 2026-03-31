# ── Build stage ───────────────────────────────────────────────────────────────
FROM node:22-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json* ./
RUN npm ci

COPY . .

# Build static export
RUN npm run build

# ── Runtime stage ─────────────────────────────────────────────────────────────
# Static export is served via a simple HTTP server
FROM node:22-alpine

RUN npm install -g serve

WORKDIR /app

COPY --from=builder /app/out ./out

EXPOSE 3000

CMD ["serve", "-s", "out", "-l", "3000"]
