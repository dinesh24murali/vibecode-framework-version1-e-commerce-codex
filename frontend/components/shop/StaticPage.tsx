export function StaticPage({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <main className="mx-auto max-w-4xl p-4">
      <h1 className="mb-4 text-3xl font-semibold">{title}</h1>
      <div className="space-y-4 text-sm leading-relaxed">{children}</div>
    </main>
  )
}

export function Section({ heading, children }: { heading: string; children: React.ReactNode }) {
  return (
    <section>
      <h2 className="mb-1 text-lg font-semibold">{heading}</h2>
      <p>{children}</p>
    </section>
  )
}

/** Inline SVG placeholder so pages need no external image hosting. */
export function PlaceholderImage({ label }: { label: string }) {
  return (
    <div role="img" aria-label={label} className="flex aspect-[4/3] items-center justify-center rounded bg-muted text-sm text-muted-foreground">
      {label}
    </div>
  )
}
