import { Button } from '@/components/ui/button'

export default function Home() {
  return (
    <main className="flex min-h-screen flex-col items-center justify-center p-24">
      <h1 className="text-4xl font-bold mb-4">e-commerce-site</h1>
      <p className="text-muted-foreground mb-8">A simple e-commerce site for selling books</p>
      <Button>Get Started</Button>
    </main>
  )
}
