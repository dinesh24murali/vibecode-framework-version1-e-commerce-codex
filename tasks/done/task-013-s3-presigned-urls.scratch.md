# Scratch: task-013 S3 presigned URLs

## Plan
1. aws-sdk-go-v2 (config, s3) via go get in the go container.
2. services/storage_service.go: StorageService interface, S3Adapter (PresignUpload with content type + length signed, PresignRead), DisabledStorage for local dev.
3. auth.RequireRole("admin") middleware; handlers/upload_handler.go; route under /admin behind auth middleware.
4. Config: optional S3_BUCKET, AWS_REGION, CDN_BASE_URL; .env.example updated. Without a bucket the endpoint returns 503 and the API still starts.
5. Contract test with the real adapter on fake env credentials (presigning is local, no AWS call) and real JWT middleware.

## Uncertainties (resolved)
- Spec vs plan: spec returns presigned_url + public_url + expires_at, 5 minute expiry, 401 for missing token; products store the public URL. Followed the spec.
- Object key products/<uuid>.<ext>; extension from content_type, the client filename is ignored (no path injection).

## Skipped
- Wiring presigned GET into product responses (not needed with a CDN). No real S3 call is made anywhere (no credentials in tests).

## Risks
- Real S3/CORS configuration is infrastructure (TASK-035). The upload itself is untested against AWS.
