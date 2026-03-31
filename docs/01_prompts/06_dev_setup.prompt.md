# Prompt: Generate Developer Setup Guide

> **How to use:** Complete Prompts 01–05 first. Fill in every `[[PLACEHOLDER]]`, then paste this entire file into your AI tool. Save the output to `docs/02_outputs/06_dev_setup.md`. This becomes the project's contributor onboarding guide.

---

## Instructions to AI

You are a senior DevOps engineer and developer advocate. Generate a **complete developer setup guide** for the project below. A brand-new engineer joining the team should be able to follow this guide and have a fully working local development environment in under 30 minutes.

---

## Project Input

**Project name:** e-commerce-site
**Backend:** Go 1.26 / Gin Web Framework
**Frontend:** NextJS 16 Static Site Generation
**Database:** Postgres
**Auth strategy:** JWT
**Containerized:** Yes — Docker Compose
**CI/CD:** none
**Environments:** dev

---

## Required Sections

1. **Prerequisites**
   - List exact tools and minimum versions required
   - Link to installation instructions for each

2. **Repository Setup**
   - Clone instructions
   - Environment variable setup (reference `.env.example`)
   - Secrets that must be obtained and from where

3. **Local Development**
   - Step-by-step: how to start all services locally
   - Expected output when everything is healthy
   - How to run just the backend / just the frontend independently

4. **Database Setup**
   - How to run migrations
   - How to seed development data
   - How to reset the database

5. **Running Tests**
   - Unit tests
   - Integration tests
   - E2e tests (`make verify`)
   - How to check coverage

6. **Common Development Tasks**
   - How to add a new API endpoint
   - How to add a new frontend page/route
   - How to create a new database migration
   - How to create a new ADR (`make adr`)
   - How to create a new task (`make task`)

7. **Debugging**
   - How to enable debug logging
   - How to connect a debugger (backend and frontend)
   - Common errors and solutions

8. **CI/CD Overview**
   - What the pipeline does on PR
   - What the pipeline does on merge to main
   - How to check pipeline status

9. **Deployment (non-prod)**
   - How to deploy to staging manually if needed

---

## Output Format

- Use Markdown with shell code blocks (` ```bash `)
- Every command must be copy-pasteable and complete
- Call out OS-specific differences (macOS vs Linux vs Windows)
- Use callout blocks for warnings: `> **Warning:** ...`
- Target length: 1500–3000 words
