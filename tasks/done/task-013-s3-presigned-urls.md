# Task: TASK-013 — S3 presigned URLs for product images

**Status:** done
**Created:** 2026-11-17
**ADR refs:** none

## Goal
Admin endpoint that issues presigned S3 PUT URLs for product image uploads.

## Acceptance Criteria
- [x] POST /api/v1/admin/upload/image/presign returns a presigned PUT URL (valid 5 minutes per the spec; the plan said 15, spec wins)
- [x] Non-admin requests return 403; unauthenticated requests return 401 (spec lists 401/403; the plan said 403 for both)
- [~] Product responses include presigned GET URLs: not done by design. The spec has the client store the CDN `public_url` and product images return that URL. `PresignRead` exists on the adapter for a private-bucket setup.
- [x] Bucket, region, CDN base URL and credentials come from env vars (S3_BUCKET, AWS_REGION, CDN_BASE_URL, AWS_ACCESS_KEY_ID/SECRET); no hardcoded values
- [x] Contract test validates the ImagePresignResponse shape (backend/tests/contract/upload_test.go)

## Dependencies
Tasks: 002, 006 (auth middleware).
