import type { Metadata } from 'next'
import { CheckoutPageContent } from '@/components/shop/CheckoutPageContent'

export const metadata: Metadata = {
  title: 'Checkout | e-commerce-site',
  description: 'Choose a delivery address, apply a coupon and place your order.',
}

export default function CheckoutPage() {
  return <CheckoutPageContent />
}
