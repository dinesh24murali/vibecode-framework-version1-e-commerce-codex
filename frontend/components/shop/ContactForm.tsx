'use client'
import { useState } from 'react'
import { submitContact } from '@/lib/api/generated/contact/contact'
import { Button } from '@/components/ui/button'

const EMAIL_RE = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
const PHONE_RE = /^\+?[0-9][0-9 ()-]{5,18}$/

type Fields = { name: string; email: string; phone: string; subject: string; message: string }
const EMPTY: Fields = { name: '', email: '', phone: '', subject: '', message: '' }

function validate(f: Fields): Partial<Record<keyof Fields, string>> {
  const e: Partial<Record<keyof Fields, string>> = {}
  if (!f.name.trim()) e.name = 'Name is required'
  if (!f.email.trim()) e.email = 'Email is required'
  else if (!EMAIL_RE.test(f.email.trim())) e.email = 'Enter a valid email address'
  if (f.phone.trim() && !PHONE_RE.test(f.phone.trim())) e.phone = 'Enter a valid phone number'
  if (!f.subject.trim()) e.subject = 'Subject is required'
  if (!f.message.trim()) e.message = 'Message is required'
  return e
}

export function ContactForm() {
  const [fields, setFields] = useState<Fields>(EMPTY)
  const [errors, setErrors] = useState<Partial<Record<keyof Fields, string>>>({})
  const [status, setStatus] = useState<'idle' | 'sending' | 'sent' | 'failed'>('idle')

  const set = (k: keyof Fields) => (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) =>
    setFields((f) => ({ ...f, [k]: e.target.value }))

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault()
    const errs = validate(fields)
    setErrors(errs)
    if (Object.keys(errs).length > 0) return
    setStatus('sending')
    try {
      const res = await submitContact({
        name: fields.name.trim(),
        email: fields.email.trim(),
        phone: fields.phone.trim() || null,
        subject: fields.subject.trim(),
        message: fields.message.trim(),
      })
      if (res.status === 200) {
        setStatus('sent')
        setFields(EMPTY)
      } else {
        setStatus('failed')
      }
    } catch {
      setStatus('failed')
    }
  }

  const input = 'mt-1 w-full rounded border px-3 py-2 text-sm'
  const field = (k: keyof Fields, label: string, type = 'text') => (
    <div>
      <label htmlFor={`contact-${k}`} className="text-sm font-medium">{label}</label>
      <input
        id={`contact-${k}`}
        type={type}
        value={fields[k]}
        onChange={set(k)}
        aria-invalid={!!errors[k]}
        aria-describedby={errors[k] ? `contact-${k}-err` : undefined}
        className={input}
      />
      {errors[k] && <p id={`contact-${k}-err`} className="mt-1 text-xs text-red-600">{errors[k]}</p>}
    </div>
  )

  return (
    <form onSubmit={onSubmit} noValidate className="space-y-3">
      {status === 'sent' && (
        <p role="status" className="rounded bg-green-50 p-3 text-sm text-green-800">
          Thank you! Your message has been sent. We will get back to you soon.
        </p>
      )}
      {status === 'failed' && (
        <p role="alert" className="rounded bg-red-50 p-3 text-sm text-red-800">
          Sorry, we could not send your message. Please try again later.
        </p>
      )}
      {field('name', 'Name')}
      {field('email', 'Email', 'email')}
      {field('phone', 'Phone (optional)', 'tel')}
      {field('subject', 'Subject')}
      <div>
        <label htmlFor="contact-message" className="text-sm font-medium">Message</label>
        <textarea
          id="contact-message"
          rows={5}
          value={fields.message}
          onChange={set('message')}
          aria-invalid={!!errors.message}
          className={input}
        />
        {errors.message && <p className="mt-1 text-xs text-red-600">{errors.message}</p>}
      </div>
      <Button type="submit" disabled={status === 'sending'}>
        {status === 'sending' ? 'Sending...' : 'Send message'}
      </Button>
    </form>
  )
}
