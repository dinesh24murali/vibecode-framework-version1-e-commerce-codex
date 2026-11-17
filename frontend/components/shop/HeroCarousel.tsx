'use client'
import { useCallback, useEffect, useState } from 'react'
import Link from 'next/link'
import useEmblaCarousel from 'embla-carousel-react'
import Autoplay from 'embla-carousel-autoplay'
import type { CarouselSlide } from '@/lib/api/generated/model'

/** Hero carousel: auto-advances every 4s, swipeable via embla's touch handling. */
export function HeroCarousel({ slides }: { slides: CarouselSlide[] }) {
  const [emblaRef, emblaApi] = useEmblaCarousel({ loop: true }, [
    Autoplay({ delay: 4000, stopOnInteraction: false }),
  ])
  const [selected, setSelected] = useState(0)
  const [broken, setBroken] = useState<Record<string, boolean>>({})

  const onSelect = useCallback(() => {
    if (emblaApi) setSelected(emblaApi.selectedScrollSnap())
  }, [emblaApi])

  useEffect(() => {
    if (!emblaApi) return
    onSelect()
    emblaApi.on('select', onSelect)
    return () => {
      emblaApi.off('select', onSelect)
    }
  }, [emblaApi, onSelect])

  // Fallback banner when there are no slides (functional spec: landing page errors).
  if (slides.length === 0) {
    return (
      <section aria-label="Featured" data-testid="carousel-fallback" className="bg-muted p-12 text-center">
        <h1 className="text-3xl font-bold">Welcome</h1>
      </section>
    )
  }

  return (
    <section aria-roledescription="carousel" aria-label="Featured" data-testid="hero-carousel" className="relative">
      <div className="overflow-hidden" ref={emblaRef}>
        <div className="flex">
          {slides.map((s) => (
            <div key={s.id} data-testid="carousel-slide" className="relative min-w-0 flex-[0_0_100%]">
              {broken[s.id] ? (
                <div className="h-64 w-full bg-muted md:h-96" />
              ) : (
                // eslint-disable-next-line @next/next/no-img-element
                <img
                  src={s.image_url}
                  alt={s.title}
                  className="h-64 w-full object-cover md:h-96"
                  onError={() => setBroken((b) => ({ ...b, [s.id]: true }))}
                />
              )}
              <div className="absolute inset-0 flex flex-col items-center justify-center gap-2 bg-black/30 p-4 text-center text-white">
                <h2 className="text-2xl font-bold md:text-4xl">{s.title}</h2>
                {s.subtitle && <p>{s.subtitle}</p>}
                {s.cta_url && s.cta_label && (
                  <Link href={s.cta_url} className="mt-2 rounded-md bg-white px-4 py-2 text-sm font-medium text-black">
                    {s.cta_label}
                  </Link>
                )}
              </div>
            </div>
          ))}
        </div>
      </div>
      <div className="absolute bottom-2 left-0 right-0 flex justify-center gap-2">
        {slides.map((s, i) => (
          <button
            key={s.id}
            aria-label={`Go to slide ${i + 1}`}
            aria-current={selected === i}
            onClick={() => emblaApi?.scrollTo(i)}
            className={`h-2 w-2 rounded-full ${selected === i ? 'bg-white' : 'bg-white/50'}`}
          />
        ))}
      </div>
    </section>
  )
}
