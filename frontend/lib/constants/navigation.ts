const slugify = (s: string) =>
  s.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/^-|-$/g, '')

export interface NavItem {
  name: string
  href: string
}

export const BOOK_CATEGORIES: NavItem[] = [
  'Colouring Books', 'Drawing Books', 'Journals', 'Calendars', 'Sketch Pads', 'Planners',
  'Stained Art', 'Fabrii Art', 'Note Book', 'Sticker Book', 'Sticker Colouring Book', 'Activity Cards',
].map((name) => ({ name, href: `/products/?category=${slugify(name)}` }))

export const PRINTING_PRODUCTS: NavItem[] = [
  'Monocarton', 'Lid and Tray Box', 'Rigid Case', 'Corrugated Box', 'Books and Pads', 'Calendars',
  'Pamplets and Brochures', 'Invitations and Business Cards', 'Envelopes and Pouches',
  'Files and Folders', 'Flyers and Danglers',
].map((name) => ({ name, href: `/printing-services/#product-${slugify(name)}` }))

export const PRINTING_SERVICES: NavItem[] = [
  'Offset Printing', 'Met. PET Printing', 'UV - Gloss/Matt Coating', 'Hot Foil Printing',
  'Screen Printing', 'Food Grade Printing', 'Rigid Case', 'Corrugation', 'Lamination', 'Pasting',
  'Die Cutting', 'Embossing/Engraving', 'Binding and Publishing',
].map((name) => ({ name, href: `/printing-services/#service-${slugify(name)}` }))

export const SIMPLE_LINKS: NavItem[] = [
  { name: 'Wall of Fame', href: '/wall-of-fame/' },
  { name: 'Contact Us', href: '/contact/' },
  { name: 'Printing Supplies', href: '/printing-supplies/' },
]
