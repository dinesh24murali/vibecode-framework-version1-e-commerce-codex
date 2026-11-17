// Business address shown on the contact page and used for the map embed.
// Configure per environment; the default is a placeholder.
export const BUSINESS_ADDRESS =
  process.env.NEXT_PUBLIC_BUSINESS_ADDRESS || '123 Printing Press Road, Chennai, Tamil Nadu 600001, India'

export const MAP_EMBED_URL = `https://www.google.com/maps?q=${encodeURIComponent(BUSINESS_ADDRESS)}&output=embed`
