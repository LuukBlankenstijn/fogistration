import { createFileRoute, redirect } from '@tanstack/react-router'

export const Route = createFileRoute('/')({
  beforeLoad: () => {
    throw redirect({
      to: "/contests"
    })
  },
  component: RouteComponent,
})

function RouteComponent() {
  return <div></div>
}
