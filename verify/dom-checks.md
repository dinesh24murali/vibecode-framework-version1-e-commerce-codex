# DOM Verification Checklist

Use this checklist after any frontend change. Work through the relevant sections depending on what changed.

---

## General (run after every frontend change)

- [ ] Page loads without console errors (`console.error` count = 0)
- [ ] Page loads without console warnings related to your change
- [ ] No broken network requests (4xx/5xx in the Network tab)
- [ ] No layout shifts visible on load (check CLS in Performance tab)
- [ ] Page title is correct
- [ ] Favicon loads

## Navigation & Routing

- [ ] All navigation links resolve to the correct routes
- [ ] Back button works correctly
- [ ] Direct URL navigation (deep links) works
- [ ] 404 page renders for unknown routes
- [ ] Redirect after login goes to the intended page

## Forms

- [ ] All required field validation triggers on submit
- [ ] Error messages are visible and descriptive
- [ ] Success state renders after valid submission
- [ ] Form clears / resets after successful submission (if expected)
- [ ] Disabled state on submit button while request is in flight
- [ ] Re-enable after response (success or error)

## Authentication

- [ ] Unauthenticated users are redirected from protected routes
- [ ] Auth token is not exposed in the DOM or URL params
- [ ] Logout clears session and redirects to login
- [ ] Expired session redirects gracefully (no blank screen)

## API Integration

- [ ] Loading states are shown while data is fetching
- [ ] Empty states render correctly when there's no data
- [ ] Error states render with a user-friendly message on API failure
- [ ] Data updates in the UI after create/update/delete operations

## Accessibility

- [ ] Focusable elements are reachable by keyboard
- [ ] Modals trap focus correctly
- [ ] `alt` text is present on all meaningful images
- [ ] Color contrast passes WCAG AA (check with browser DevTools)
- [ ] `aria-label` or visible labels present on all form inputs

## Responsive Design

- [ ] Mobile viewport (375px): no horizontal overflow
- [ ] Tablet viewport (768px): layout is usable
- [ ] Desktop viewport (1280px+): no stretched or misaligned elements

## Performance

- [ ] No synchronous blocking scripts in `<head>`
- [ ] Images are appropriately sized (no massive images for small displays)
- [ ] Largest Contentful Paint (LCP) is under 2.5s on a throttled connection
