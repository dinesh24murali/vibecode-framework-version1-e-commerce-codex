'use client'
import { useState } from 'react'
import { useRouter } from 'next/navigation'

export function SearchBar({ onSubmit }: { onSubmit?: () => void }) {
  const router = useRouter()
  const [value, setValue] = useState('')

  return (
    <form
      role="search"
      onSubmit={(e) => {
        e.preventDefault()
        const q = value.trim()
        router.push(q ? `/products/?q=${encodeURIComponent(q)}` : '/products/')
        onSubmit?.()
      }}
    >
      <input
        type="search"
        aria-label="Search store"
        placeholder="Search books..."
        value={value}
        onChange={(e) => setValue(e.target.value)}
        className="h-10 w-full rounded border px-3 text-sm"
      />
    </form>
  )
}
