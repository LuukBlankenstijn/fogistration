import { UserRole } from '@/clients/generated-client'
import { createFileRoute, Outlet, redirect } from '@tanstack/react-router'

export const Route = createFileRoute('/_private/admin')({
  component: RouteComponent,
  beforeLoad: ({ context }) => {
    if (context.auth.user?.role !== UserRole.ADMIN) {
      throw redirect({ to: "/dashboard" })
    }
  }
})

function RouteComponent() {
  return <Outlet />
}
