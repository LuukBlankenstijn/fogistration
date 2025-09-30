import TeamsPage from '@/pages/TeamsPage'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_sidebar/teams')({
  component: RouteComponent,
})

function RouteComponent() {
  return <TeamsPage />
}
