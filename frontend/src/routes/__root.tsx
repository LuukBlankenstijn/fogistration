import { Outlet, createRootRoute } from '@tanstack/react-router'
import { TanStackRouterDevtools } from '@tanstack/react-router-devtools'


export const Route = createRootRoute({
  component: () => (
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
      <div className='relative z-10'>
        <Outlet />
        <TanStackRouterDevtools />
      </div>
    </div>
  ),
})
