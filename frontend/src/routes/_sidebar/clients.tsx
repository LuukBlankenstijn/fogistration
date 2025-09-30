import ClientsPage from '@/pages/ClientsPage'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_sidebar/clients')({
  component: RouteComponent,
})

function RouteComponent() {
  return <ClientsPage />
}
