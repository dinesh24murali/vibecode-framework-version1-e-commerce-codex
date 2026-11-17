'use client'
import { useEffect, useState } from 'react'
import {
  changePassword,
  confirmEmailChange,
  getProfile,
  requestEmailChange,
  updateProfile,
} from '@/lib/api/generated/account/account'
import type { UserProfile } from '@/lib/api/generated/model'
import { errorMessage } from '@/lib/auth/helpers'
import { useStore } from '@/lib/store'
import { Button } from '@/components/ui/button'

const input = 'w-full rounded-md border px-3 py-2'

function Field({
  id,
  label,
  ...rest
}: { id: string; label: string } & React.InputHTMLAttributes<HTMLInputElement>) {
  return (
    <div className="space-y-1">
      <label htmlFor={id} className="block text-sm font-medium">
        {label}
      </label>
      <input id={id} className={input} {...rest} />
    </div>
  )
}

export function ProfileTab() {
  const setUser = useStore((s) => s.setUser)
  const [profile, setProfile] = useState<UserProfile | null>(null)
  const [first, setFirst] = useState('')
  const [last, setLast] = useState('')
  const [msg, setMsg] = useState<{ ok: boolean; text: string } | null>(null)

  const [current, setCurrent] = useState('')
  const [newPw, setNewPw] = useState('')
  const [confirmPw, setConfirmPw] = useState('')

  const [newEmail, setNewEmail] = useState('')
  const [otp, setOtp] = useState('')
  const [otpSent, setOtpSent] = useState(false)

  function apply(p: UserProfile) {
    setProfile(p)
    setFirst(p.first_name)
    setLast(p.last_name)
    setUser({ id: p.id, email: p.email, first_name: p.first_name, last_name: p.last_name, role: p.role })
  }

  useEffect(() => {
    void getProfile().then((res) => {
      if (res.status === 200) apply(res.data as UserProfile)
      else setMsg({ ok: false, text: errorMessage(res) })
    })
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  async function saveProfile(e: React.FormEvent) {
    e.preventDefault()
    const res = await updateProfile({ first_name: first, last_name: last })
    if (res.status === 200) {
      apply(res.data as UserProfile)
      setMsg({ ok: true, text: 'Profile saved' })
    } else setMsg({ ok: false, text: errorMessage(res) })
  }

  async function savePassword(e: React.FormEvent) {
    e.preventDefault()
    const res = await changePassword({
      current_password: current,
      new_password: newPw,
      confirm_new_password: confirmPw,
    })
    if (res.status === 200) {
      setCurrent('')
      setNewPw('')
      setConfirmPw('')
      setMsg({ ok: true, text: 'Password updated' })
    } else setMsg({ ok: false, text: errorMessage(res) })
  }

  async function sendOtp(e: React.FormEvent) {
    e.preventDefault()
    const res = await requestEmailChange({ new_email: newEmail })
    if (res.status === 200) {
      setOtpSent(true)
      setMsg({ ok: true, text: `Verification code sent to ${newEmail}` })
    } else setMsg({ ok: false, text: errorMessage(res) })
  }

  async function confirmEmail(e: React.FormEvent) {
    e.preventDefault()
    const res = await confirmEmailChange({ new_email: newEmail, otp })
    if (res.status === 200) {
      apply(res.data as UserProfile)
      setOtpSent(false)
      setNewEmail('')
      setOtp('')
      setMsg({ ok: true, text: 'Email updated' })
    } else setMsg({ ok: false, text: errorMessage(res) })
  }

  return (
    <div className="space-y-8">
      {msg && (
        <p role="status" className={msg.ok ? 'text-sm text-green-700' : 'text-sm text-red-600'}>
          {msg.text}
        </p>
      )}
      <form onSubmit={saveProfile} className="space-y-3">
        <h2 className="text-lg font-medium">Profile</h2>
        <Field id="profile-first" label="First name" required value={first} onChange={(e) => setFirst(e.target.value)} />
        <Field id="profile-last" label="Last name" required value={last} onChange={(e) => setLast(e.target.value)} />
        <Button type="submit">Save profile</Button>
      </form>

      <form onSubmit={otpSent ? confirmEmail : sendOtp} className="space-y-3">
        <h2 className="text-lg font-medium">Email</h2>
        <p className="text-sm">Current: {profile?.email}</p>
        <Field id="new-email" label="New email" type="email" required value={newEmail} onChange={(e) => setNewEmail(e.target.value)} disabled={otpSent} />
        {otpSent && (
          <Field id="email-otp" label="Verification code" inputMode="numeric" maxLength={6} required value={otp} onChange={(e) => setOtp(e.target.value)} />
        )}
        <Button type="submit">{otpSent ? 'Confirm email change' : 'Send verification code'}</Button>
      </form>

      <form onSubmit={savePassword} className="space-y-3">
        <h2 className="text-lg font-medium">Change password</h2>
        <Field id="pw-current" label="Current password" type="password" required value={current} onChange={(e) => setCurrent(e.target.value)} />
        <Field id="pw-new" label="New password" type="password" required minLength={8} value={newPw} onChange={(e) => setNewPw(e.target.value)} />
        <Field id="pw-confirm" label="Confirm new password" type="password" required value={confirmPw} onChange={(e) => setConfirmPw(e.target.value)} />
        <Button type="submit">Update password</Button>
      </form>
    </div>
  )
}
