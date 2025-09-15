import { getActiveContestOptions } from '@/clients/generated-client/@tanstack/react-query.gen'
import { ClientsTable } from '@/components/ClientsTable'
import { ContestTable } from '@/components/ContestTable'
import { Section } from '@/components/Section'
import { TeamsTable } from '@/components/TeamsTable'
import queryClient from '@/query/query-client'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_private/dashboard')({
  loader: () => {
    void queryClient.ensureQueryData(getActiveContestOptions())
  },
  component: Dashboard,
})


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

