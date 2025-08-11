import { useAuth } from '@/auth'
import { createFileRoute, Outlet, redirect, useNavigate } from '@tanstack/react-router'
import { useEffect } from 'react'

export const Route = createFileRoute('/_private')({
  beforeLoad: ({ context, location }) => {
    if (!context.auth.isAuthenticated) {
      // eslint-disable-next-line @typescript-eslint/only-throw-error
      throw redirect({
        to: "/",
        search: {
          redirect: location.href
        }
      })
    }
  },
  component: PrivateComponent,
})

function PrivateComponent() {
  const { isAuthenticated } = useAuth()
  const navigate = useNavigate()

  useEffect(() => {
    if (!isAuthenticated) {
      void navigate({ to: "/" })
    }
  }, [isAuthenticated])

  return <Outlet />
}
