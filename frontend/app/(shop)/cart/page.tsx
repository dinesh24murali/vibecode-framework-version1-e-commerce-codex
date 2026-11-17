import type { Metadata } from 'next'
import { CartPageContent } from '@/components/shop/CartPageContent'

export const metadata: Metadata = {
  title: 'Your Cart | e-commerce-site',
  description: 'Review the items in your cart and proceed to checkout.',
}

export default function CartPage() {
  return <CartPageContent />
}
