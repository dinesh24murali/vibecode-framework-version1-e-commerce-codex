# Scratch: add-forgot-password-flow

## Plan

1. **API spec** (`04_api_spec.yaml`):
   a. Insert 3 new paths after `/api/v1/auth/logout` and before the `# Landing` section
   b. Insert 2 new schemas (`PasswordResetRequestRequest`, `PasswordResetCompleteRequest`) after `RegisterCompleteRequest`

2. **Functional spec** (`02_functional_spec.md`):
   a. Update `/login` row in Screen Inventory — remove "if added later" caveat
   b. Add UF-015 after UF-014
   c. Expand FS-006 with forgot password happy path, error states, and validation rules
   d. Add password reset OTP row to Notifications & Emails table

## Files to Change

- `docs/02_outputs/04_api_spec.yaml`
- `docs/02_outputs/02_functional_spec.md`

## Uncertainties

None — insertion points confirmed by reading the files.

## What Will Be Skipped

- `docs/02_outputs/03_tech_architecture.md` — no changes needed
- No `/forgot-password` route added to screen inventory (it's a sub-flow of `/login`, handled by frontend routing)

## Risks

- `VerifyOTPRequest` is currently used for registration only; repurposing it for password reset is valid since the shape is identical (`{email, otp}`). No changes needed.
- `OTPTokenResponse` similarly reused as-is.
