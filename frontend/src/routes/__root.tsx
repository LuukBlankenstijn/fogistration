import { Outlet, createRootRoute } from '@tanstack/react-router'
import { TanStackRouterDevtools } from '@tanstack/react-router-devtools'


function RootLayout() {
  return (
    <div className="relative min-h-dvh overflow-hidden bg-[hsl(var(--bg))]">
      <div aria-hidden className="pointer-events-none absolute inset-0">
        <div className="absolute -top-32 -left-32 h-96 w-96 rounded-full bg-[hsla(var(--brand)/0.12)] blur-3xl" />
        <div className="absolute -bottom-32 -right-32 h-96 w-96 rounded-full bg-[hsla(var(--brand)/0.10)] blur-3xl" />
      </div>
      <div
        aria-hidden
        className="pointer-events-none absolute inset-0 opacity-40"
        style={{
          backgroundImage:
            'linear-gradient(to right, hsl(var(--border)) 1px, transparent 1px), linear-gradient(to bottom, hsl(var(--border)) 1px, transparent 1px)',
          backgroundSize: '24px 24px'
        }}
      />

      <div className="relative flex min-h-dvh">
        <Outlet />
        <TanStackRouterDevtools />
      </div>
    </div>
  )
}

export const Route = createRootRoute({
  component: () => (
    <RootLayout />
  ),
})
