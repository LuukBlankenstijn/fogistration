import { useAuth } from '@/auth';
import { LogoutButton } from '@/components/logout';
import { createFileRoute } from '@tanstack/react-router'
import { useMemo } from 'react';

export const Route = createFileRoute('/_private/dashboard')({
  component: RouteComponent,
})

function RouteComponent() {
  return <Dashboard />
}

function StatCard(props: { label: string; value: string; sub?: string }) {
  return (
    <div className="rounded-2xl border border-[hsl(var(--border))] bg-[hsl(var(--panel))] p-4 shadow-sm">
      <div className="text-sm text-[hsl(var(--muted))]">{props.label}</div>
      <div className="mt-1 text-2xl font-semibold text-[hsl(var(--fg))]">{props.value}</div>
      {props.sub && <div className="mt-1 text-xs text-[hsl(var(--muted))]">{props.sub}</div>}
    </div>
  )
}

function QuickAction(props: React.ButtonHTMLAttributes<HTMLButtonElement> & { icon?: React.ReactNode }) {
  const { className = '', children, ...rest } = props
  return (
    <button
      type="button"
      {...rest}
      className={`rounded-lg border border-[hsl(var(--border))] bg-[hsl(var(--input))] px-3 py-2 text-sm transition hover:bg-[hsl(var(--hover))] ${className}`}
    >
      {children}
    </button>
  )
}

function Section(props: { title: string; children: React.ReactNode; right?: React.ReactNode }) {
  return (
    <section className="rounded-2xl border border-[hsl(var(--border))] bg-[hsl(var(--panel))] p-4 shadow-sm">
      <header className="mb-3 flex items-center justify-between">
        <h2 className="text-sm font-semibold text-[hsl(var(--fg))]">{props.title}</h2>
        {props.right}
      </header>
      {props.children}
    </section>
  )
}

function Table(props: { headers: string[]; rows: React.ReactNode[][] }) {
  return (
    <div className="overflow-x-auto">
      <table className="w-full border-collapse text-sm">
        <thead>
          <tr className="text-left text-[hsl(var(--muted))]">
            {props.headers.map((h) => (
              <th key={h} className="border-b border-[hsl(var(--border))] px-3 py-2 font-medium">
                {h}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {props.rows.map((r, i) => (
            <tr key={i} className="hover:bg-[hsl(var(--hover))]">
              {r.map((c, j) => (
                <td key={j} className="border-b border-[hsl(var(--border))] px-3 py-2">
                  {c}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

function Dashboard() {
  const { user, logout } = useAuth()

  const stats = useMemo(
    () => [
      { label: 'Projects', value: '12' },
      { label: 'Active Sessions', value: '3' },
      { label: 'Builds Today', value: '7', sub: '+2 since yesterday' },
      { label: 'Errors', value: '0', sub: 'Last 24h' },
    ],
    []
  )

  const recent = useMemo(
    () => [
      ['#1432', 'Deploy', 'api-gateway', '2m ago'],
      ['#1431', 'Job', 'daily-backup', '1h ago'],
      ['#1430', 'Deploy', 'web-frontend', '3h ago'],
      ['#1429', 'Alert', 'db-latency', 'Yesterday'],
    ],
    []
  )

  return (
    <div className="min-h-screen bg-[hsl(var(--bg))] text-[hsl(var(--fg))]">
      <div className="mx-auto max-w-6xl p-6">
        <header className="mb-6">
          <h1 className="text-xl font-bold text-[hsl(var(--brand))]">Fogistration</h1>
          <div className="mt-1 text-lg font-semibold">
            Welcome{user?.username ? `, ${user.username}` : ''} 👋
          </div>
          <LogoutButton onClick={logout} />
        </header>

        {/* Quick actions */}
        <div className="mb-6 grid grid-cols-1 gap-3 sm:grid-cols-2 md:grid-cols-4">
          <QuickAction onClick={() => { /* no-op */ }}>New Project</QuickAction>
          <QuickAction onClick={() => { /* no-op */ }}>Invite User</QuickAction>
          <QuickAction onClick={() => { /* no-op */ }}>Run Build</QuickAction>
          <QuickAction onClick={() => { /* no-op */ }}>View Logs</QuickAction>
        </div>

        {/* Stats */}
        <div className="mb-6 grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
          {stats.map((s) => (
            <StatCard key={s.label} label={s.label} value={s.value} sub={s.sub} />
          ))}
        </div>

        <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
          {/* Activity */}
          <div className="lg:col-span-2">
            <Section
              title="Recent Activity"
              right={
                <button className="text-xs text-[hsl(var(--brand))] hover:underline">View all</button>
              }
            >
              <Table
                headers={['ID', 'Type', 'Target', 'When']}
                rows={recent.map((r) => r.map((c) => <span key={c}>{c}</span>))}
              />
            </Section>
          </div>

          {/* Announcements / Tips */}
          <div className="lg:col-span-1">
            <Section title="Tips">
              <ul className="space-y-3 text-sm">
                <li className="rounded-lg bg-[hsl(var(--input))] p-3">
                  Use access tokens with short TTL and refresh tokens for longer sessions.
                </li>
                <li className="rounded-lg bg-[hsl(var(--input))] p-3">
                  Group routes in Huma to share auth & role middleware.
                </li>
                <li className="rounded-lg bg-[hsl(var(--input))] p-3">
                  Prefetch route data with TanStack Router loaders for instant views.
                </li>
              </ul>
            </Section>
          </div>
        </div>
      </div>
    </div>
  )
}
