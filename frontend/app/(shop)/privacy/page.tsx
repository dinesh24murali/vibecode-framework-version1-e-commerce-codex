import type { Metadata } from 'next'
import { Section, StaticPage } from '@/components/shop/StaticPage'

export const metadata: Metadata = {
  title: 'Privacy Policy | e-commerce-site',
  description: 'How we collect, use and protect your personal information.',
}

export default function PrivacyPage() {
  return (
    <StaticPage title="Privacy Policy">
      <p>Boilerplate text only; have it reviewed by counsel before launch.</p>
      <Section heading="1. Information We Collect">
        Account details, delivery addresses, order history and messages you send us.
      </Section>
      <Section heading="2. How We Use Information">
        To process orders, provide support, send order emails and improve the site.
      </Section>
      <Section heading="3. Sharing">We do not sell personal data. We share it only with service providers needed to fulfil orders.</Section>
      <Section heading="4. Cookies">We use cookies for sign-in sessions and basic site functionality.</Section>
      <Section heading="5. Data Retention and Security">
        We keep data only as long as needed and protect it with industry-standard safeguards.
      </Section>
      <Section heading="6. Your Rights">You may request access, correction or deletion of your data through the Contact Us page.</Section>
    </StaticPage>
  )
}
