import Sidebar from '@/components/SideBar'
import { createFileRoute, Outlet } from '@tanstack/react-router'

export const Route = createFileRoute('/_sidebar')({
  component: RouteComponent,
})

function RouteComponent() {
  return (
    <>
      <Sidebar />
      <main className="flex-1 p-6 space-y-4 overflow-y-auto">
        <Outlet />
      </main>
    </ >
  )
}
