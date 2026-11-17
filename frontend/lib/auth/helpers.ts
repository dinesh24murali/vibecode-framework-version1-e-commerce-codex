/** Extracts a human-readable message from a generated-client response. */
export function errorMessage(res: { data: unknown }, fallback = 'Something went wrong'): string {
  const d = res.data as { message?: string } | undefined
  return d?.message || fallback
}

/** Only allow same-origin relative paths to avoid open redirects. */
export function safeRedirect(target: string | null): string {
  if (target && target.startsWith('/') && !target.startsWith('//') && !target.startsWith('/\\')) {
    return target
  }
  return '/'
}
