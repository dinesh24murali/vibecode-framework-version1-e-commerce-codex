'use client'
import { useState } from 'react'
import { validateCoupon } from '@/lib/api/generated/checkout/checkout'
import type { CouponValidationResponse } from '@/lib/api/generated/model'

interface Props {
  /** Called with the code when it is valid, or null when empty/invalid. */
  onApply: (code: string | null) => void
}

export function CouponInput({ onApply }: Props) {
  const [value, setValue] = useState('')
  const [feedback, setFeedback] = useState<{ valid: boolean; message: string } | null>(null)

  async function check() {
    const code = value.trim()
    if (!code) {
      setFeedback(null)
      onApply(null)
      return
    }
    try {
      const res = await validateCoupon({ coupon_code: code })
      if (res.status !== 200) throw new Error('validate failed')
      const body = res.data as CouponValidationResponse
      setFeedback({ valid: body.is_valid, message: body.message })
      onApply(body.is_valid ? body.coupon_code : null)
    } catch {
      setFeedback({ valid: false, message: 'Could not check the coupon. Please try again.' })
      onApply(null)
    }
  }

  return (
    <div>
      <label htmlFor="coupon" className="mb-1 block text-sm font-medium">Coupon code</label>
      <input
        id="coupon"
        value={value}
        onChange={(e) => setValue(e.target.value)}
        onBlur={check}
        className="w-full rounded border px-3 py-2 uppercase"
        autoComplete="off"
      />
      {feedback && (
        <p data-testid="coupon-feedback" className={`mt-1 text-sm ${feedback.valid ? 'text-green-700' : 'text-red-700'}`}>
          {feedback.message}
        </p>
      )}
    </div>
  )
}
