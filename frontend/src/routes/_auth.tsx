import { useAuth } from '@/auth'
import { createFileRoute, Outlet, redirect, useNavigate } from '@tanstack/react-router'
import { useEffect } from 'react'

export const Route = createFileRoute('/_auth')({
  validateSearch: (search) => ({
    redirect: (search.redirect as string) || "/"
  }),
  beforeLoad: ({ context, location }) => {
    if (context.auth.isAuthenticated) {
      throw redirect({
        to: "/",
        search: {
          redirect: location.href
        }
      })
    }
  },
  component: PublicComponent,
})

function PublicComponent() {
  const { isAuthenticated } = useAuth()
  const { redirect } = Route.useSearch()
  const navigate = useNavigate()


  useEffect(() => {
    if (isAuthenticated) {
      void navigate({ to: redirect })
    }
  }, [isAuthenticated])

  return <Outlet />
}
