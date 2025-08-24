import { getActiveContestOptions, getAllClientsOptions, getAllTeamsOptions } from '@/clients/generated-client/@tanstack/react-query.gen'
import { Card } from '@/components/Card'
import { ClientsTable } from '@/components/ClientsTable'
import { ContestTable } from '@/components/ContestTable'
import { TeamsTable } from '@/components/TeamsTable'
import queryClient from '@/query/client'
import { createFileRoute } from '@tanstack/react-router'
import { Suspense, type ReactNode } from 'react'

export const Route = createFileRoute('/_private/dashboard')({
  loader: () => {
    void queryClient.ensureQueryData(getAllClientsOptions())
    void queryClient.ensureQueryData(getAllTeamsOptions())
    void queryClient.ensureQueryData(getActiveContestOptions())
  },
  component: Dashboard,
})

function Section({ children, fallback }: { children: ReactNode, fallback: string }) {
  return (
    <Card className='m-6 max-w-6xl'>
      <Suspense fallback={<div className="text-sm text-[hsl(var(--muted))]">{fallback}</div>}>
        <section className='space-y-3'>
          {children}
        </section>
      </Suspense>
    </Card>
  )
}

export default function Dashboard() {
  return (
    <main className="relative mx-auto max-w-6xl px-6 py-8">
      <Section fallback='Loading clients...'>
        <ClientsTable />
      </Section>

      <Section fallback='Loading teams...'>
        <TeamsTable />
      </Section>

      <Section fallback='Loading contests...'>
        <ContestTable />
      </Section>
    </main>
  )
}

