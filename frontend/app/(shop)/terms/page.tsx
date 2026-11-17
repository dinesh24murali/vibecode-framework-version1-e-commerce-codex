import type { Metadata } from 'next'
import { Section, StaticPage } from '@/components/shop/StaticPage'

export const metadata: Metadata = {
  title: 'Terms & Conditions | e-commerce-site',
  description: 'The terms and conditions that apply to purchases and use of this website.',
}

export default function TermsPage() {
  return (
    <StaticPage title="Terms & Conditions">
      <p>Boilerplate text only; have it reviewed by counsel before launch.</p>
      <Section heading="1. Acceptance of Terms">By using this site you agree to these terms.</Section>
      <Section heading="2. Orders and Pricing">
        Prices are in INR and include applicable taxes unless stated. We may cancel orders affected by pricing or stock errors.
      </Section>
      <Section heading="3. Shipping and Delivery">Delivery estimates are indicative and not guaranteed.</Section>
      <Section heading="4. Cancellations and Returns">
        Open orders may be cancelled by the customer before fulfilment. Return eligibility is described at checkout.
      </Section>
      <Section heading="5. Intellectual Property">All content on this site belongs to us or our licensors.</Section>
      <Section heading="6. Limitation of Liability">
        To the extent permitted by law, our liability is limited to the value of the order concerned.
      </Section>
      <Section heading="7. Contact">Questions about these terms can be sent through the Contact Us page.</Section>
    </StaticPage>
  )
}
