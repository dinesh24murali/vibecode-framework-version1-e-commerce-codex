import type { Metadata } from 'next'
import './globals.css'
import { StoreHydrator } from '@/components/StoreHydrator'

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
      <body>
        <StoreHydrator />
        {children}
      </body>
    </html>
  )
}
