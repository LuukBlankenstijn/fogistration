import { createFileRoute, redirect } from '@tanstack/react-router'

export const Route = createFileRoute('/')({
  beforeLoad: ({ context }) => {
    if (!context.auth.isAuthenticated) {
      // eslint-disable-next-line @typescript-eslint/only-throw-error
      throw redirect({
        to: "/login",
        search: {
          redirect: "/dashboard"
        }
      })
    } else {
      // eslint-disable-next-line @typescript-eslint/only-throw-error
      throw redirect({
        to: "/dashboard"
      })
    }
  },
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/"!</div>
}
