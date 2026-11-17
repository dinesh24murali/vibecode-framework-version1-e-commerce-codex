# Scratch: task-009 auth pages

## Plan
1. lib/auth/helpers.ts: errorMessage(), safeRedirect() (same-origin paths only).
2. components/auth: OtpInput, LoginForm, SignupForm (step 1), VerifyOtpForm (step 2 OTP then step 3 name+password on the same page, otp_token in state), FacebookLoginButton (FB SDK script, popup, loginFacebook).
3. Pages (auth)/login, signup, signup/verify-otp; useSearchParams components wrapped in Suspense (static export).
4. Tokens -> authSlice.login(); redirect via ?redirect=.
5. e2e: Playwright with page.route mocking /api/v1/auth/* (the api service needs JWT keys and is not part of the static verify stack). Documented decision.
6. Add NEXT_PUBLIC_FACEBOOK_APP_ID to .env.example.

## Files
See task file + frontend/lib/auth/helpers.ts, .env.example

## Uncertainties (resolved)
- Step 3 placed on verify-otp page (plan lists only 3 pages). Email passed via ?email=.
- "against running backend": backend covered by contract tests; e2e mocks network.
- Profile fetch after login skipped (no profile endpoint in scope); user stays null until account tasks.

## Skipped
- Forgot-password UI (task 010).

## Risks
- Playwright config dir mode (tests/e2e) must pick up the spec.
