# Task: add-forgot-password-flow

**Status:** active
**Created:** 2026-03-26
**ADR refs:** none

---

## Goal

Add the forgot password flow to the API spec and functional spec so customers can reset their password via an OTP-verified 3-step flow.

## Background

The `/login` page in the functional spec noted "forgot password link if added later" — this was intentionally deferred. It is now in scope. The tech architecture already supports OTP flows (same pattern as registration), so only the spec documents need updating. The existing `VerifyOTPRequest` and `OTPTokenResponse` schemas are reused directly.

## Acceptance Criteria

- [ ] 3 new auth endpoints added to API spec: `password-reset/request`, `password-reset/verify-otp`, `password-reset/complete`
- [ ] 2 new schemas added: `PasswordResetRequestRequest`, `PasswordResetCompleteRequest`
- [ ] Existing schemas reused (not duplicated): `VerifyOTPRequest`, `OTPTokenResponse`, `MessageResponse`
- [ ] UF-015 (Forgot Password) added to functional spec
- [ ] FS-006 expanded with forgot password happy path, error states, and validation rules
- [ ] Screen inventory `/login` entry updated (remove "if added later" caveat)
- [ ] Notifications table updated with password reset OTP row

## Dependencies

- **Tasks:** none
- **Memory files:** none
- **Docs:** `docs/02_outputs/04_api_spec.yaml`, `docs/02_outputs/02_functional_spec.md`

## Files Expected to Change

- `docs/02_outputs/04_api_spec.yaml`
- `docs/02_outputs/02_functional_spec.md`

## Notes

- `password-reset/request` must always respond 200 even when the email is not found (prevents account enumeration)
- OTP expiry is 10 minutes — same as registration OTP
- Tech architecture doc needs no changes

---

## Scratch File

See `tasks/active/add-forgot-password-flow.scratch.md`
