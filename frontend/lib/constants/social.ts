// Social handles are configured per environment; defaults are the platform home pages.
export const SOCIAL_LINKS = [
  { name: 'Facebook', href: process.env.NEXT_PUBLIC_FACEBOOK_URL || 'https://www.facebook.com/' },
  { name: 'Instagram', href: process.env.NEXT_PUBLIC_INSTAGRAM_URL || 'https://www.instagram.com/' },
  { name: 'WhatsApp', href: process.env.NEXT_PUBLIC_WHATSAPP_URL || 'https://wa.me/' },
] as const
