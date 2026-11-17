'use client'

interface OtpInputProps {
  value: string
  onChange: (v: string) => void
  error?: string | null
}

export function OtpInput({ value, onChange, error }: OtpInputProps) {
  return (
    <div className="space-y-1">
      <label htmlFor="otp" className="block text-sm font-medium">
        6-digit code
      </label>
      <input
        id="otp"
        name="otp"
        inputMode="numeric"
        autoComplete="one-time-code"
        maxLength={6}
        value={value}
        onChange={(e) => onChange(e.target.value.replace(/\D/g, '').slice(0, 6))}
        aria-invalid={error ? true : undefined}
        className="w-full rounded-md border px-3 py-2 text-center text-lg tracking-widest"
      />
      {error && (
        <p role="alert" className="text-sm text-red-600">
          {error}
        </p>
      )}
    </div>
  )
}
