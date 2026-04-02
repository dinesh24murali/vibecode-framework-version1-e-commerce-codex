import type { Metadata } from 'next'
import './globals.css'

export const metadata: Metadata = {
  title: 'e-commerce-site',
  description: 'A simple e-commerce site for selling books',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  )
}
