# Task: TASK-019 — Static pages

**Status:** done
**Created:** 2026-11-17
**ADR refs:** none

## Goal
Wall of Fame, About Us, Terms & Conditions and Privacy Policy as static pages.

## Acceptance Criteria
- [x] All four routes exist in the static export (e2e loads each; build lists them)
- [x] Each page has a unique, descriptive title
- [x] No client-side JS errors on any of the four pages (e2e asserts)
- [x] Static build, e2e and DOM checks pass (docker equivalents of make verify)

## Dependencies
Tasks: 003.
