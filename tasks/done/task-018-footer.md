# Task: TASK-018 — Footer component

**Status:** done
**Created:** 2026-11-17
**ADR refs:** none

## Goal
Site footer with internal links and social icons on every shop page.

## Acceptance Criteria
- [x] Internal links (Terms & Conditions, About Us, Wall of Fame) point to /terms/, /about/, /wall-of-fame/ (pages arrive in TASK-019)
- [x] Social icons open in a new tab with rel noopener noreferrer
- [x] Footer present on every page via the (shop) layout
- [x] Static build, e2e and DOM checks pass (docker equivalents of make verify)

## Dependencies
Tasks: 003. Memory: memory/product.md.
