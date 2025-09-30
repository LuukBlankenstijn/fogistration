import ContestPage from '@/pages/ContestPage'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_sidebar/contests')({
  component: RouteComponent,
})

function RouteComponent() {
  return <ContestPage />
}
