import { getActiveContestOptions, getAllClientsOptions, getAllTeamsOptions } from '@/clients/generated-client/@tanstack/react-query.gen'
import { ClientsTable } from '@/components/ClientsTable'
import { TeamsTable } from '@/components/TeamsTable'
import queryClient from '@/query/client'
import { createFileRoute } from '@tanstack/react-router'
import { Suspense } from 'react'

export const Route = createFileRoute('/_private/dashboard')({
  loader: () => {
    void queryClient.ensureQueryData(getAllClientsOptions())
    void queryClient.ensureQueryData(getAllTeamsOptions())
    void queryClient.ensureQueryData(getActiveContestOptions())
  },
  component: Dashboard,
})

function ClientsSection() {
  return (
    <section className="space-y-3">
      <ClientsTable />
    </section>
  )
}

function TeamsSection() {
  return (
    <section className="space-y-3">
      <TeamsTable />
    </section>
  )
}

export default function Dashboard() {
  return (
    <div className="relative min-h-screen overflow-hidden bg-[hsl(var(--bg))]">
      {/* subtle brand glow */}
      <div aria-hidden className="pointer-events-none absolute inset-0">
        <div className="absolute -top-32 -left-32 h-96 w-96 rounded-full bg-[hsla(var(--brand)/0.12)] blur-3xl" />
        <div className="absolute -bottom-32 -right-32 h-96 w-96 rounded-full bg-[hsla(var(--brand)/0.10)] blur-3xl" />
      </div>

      {/* faint grid */}
      <div
        aria-hidden
        className="pointer-events-none absolute inset-0 opacity-40"
        style={{
          backgroundImage:
            'linear-gradient(to right, hsl(var(--border)) 1px, transparent 1px), linear-gradient(to bottom, hsl(var(--border)) 1px, transparent 1px)',
          backgroundSize: '24px 24px'
        }}
      />

      <main className="relative mx-auto max-w-6xl px-6 py-8">
        <section className="rounded-2xl border border-[hsl(var(--border))] bg-[hsl(var(--panel))] p-6 m-6 shadow-sm">
          <header className="mb-6 flex items-center justify-between">
            {/* <h1 className="text-2xl font-bold">Dashboard</h1> */}
            {/* toolbar later */}
          </header>

          <Suspense fallback={<div className="text-sm text-[hsl(var(--muted))]">Loading clients…</div>}>
            <ClientsSection />
          </Suspense>
        </section>


        <section className="rounded-2xl border border-[hsl(var(--border))] bg-[hsl(var(--panel))] p-6 m-6 shadow-sm">
          <header className="mb-6 flex items-center justify-between">
            {/* <h1 className="text-2xl font-bold">Dashboard</h1> */}
            {/* toolbar later */}
          </header>

          <Suspense fallback={<div className="text-sm text-[hsl(var(--muted))]">Loading teams…</div>}>
            <TeamsSection />
          </Suspense>
        </section>
      </main>
    </div>
  )
}

