'use client'
import { useRef, useState } from 'react'

const dist = (t: React.TouchList) => Math.hypot(t[0].clientX - t[1].clientX, t[0].clientY - t[1].clientY)

/** Hover-zoom on desktop (cursor-follow), two-finger pinch zoom on touch devices. */
export function ProductZoom({ src, alt }: { src: string; alt: string }) {
  const [origin, setOrigin] = useState('50% 50%')
  const [scale, setScale] = useState(1)
  const pinch = useRef<{ d: number; s: number } | null>(null)

  return (
    <div
      data-testid="product-zoom"
      className="relative aspect-square w-full overflow-hidden rounded border"
      style={{ touchAction: 'pan-y' }}
      onMouseEnter={() => setScale(2)}
      onMouseLeave={() => {
        setScale(1)
        setOrigin('50% 50%')
      }}
      onMouseMove={(e) => {
        const r = e.currentTarget.getBoundingClientRect()
        setOrigin(`${((e.clientX - r.left) / r.width) * 100}% ${((e.clientY - r.top) / r.height) * 100}%`)
      }}
      onTouchStart={(e) => {
        if (e.touches.length === 2) pinch.current = { d: dist(e.touches), s: scale }
      }}
      onTouchMove={(e) => {
        if (e.touches.length === 2 && pinch.current) {
          setScale(Math.min(4, Math.max(1, (pinch.current.s * dist(e.touches)) / pinch.current.d)))
        }
      }}
      onTouchEnd={() => {
        pinch.current = null
      }}
    >
      {/* eslint-disable-next-line @next/next/no-img-element */}
      <img
        src={src}
        alt={alt}
        draggable={false}
        data-testid="zoom-image"
        data-scale={scale}
        className="h-full w-full object-cover transition-transform duration-100"
        style={{ transform: `scale(${scale})`, transformOrigin: origin }}
      />
    </div>
  )
}
