import { useAuth } from '@/auth'
import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'

export const Route = createFileRoute('/_auth/login')({
  component: Login,
})

function Login() {
  const [username, setUsername] = useState("")
  const [password, setPassword] = useState("")
  const { login, devLogin } = useAuth()

  const signIn = () => {
    login(username, password)
  }

  return (
    <div className="min-h-screen grid place-items-center bg-[hsl(var(--bg))] text-[hsl(var(--fg))]">
      <div className="w-full max-w-md rounded-2xl border border-[hsl(var(--border))] bg-[hsl(var(--panel))] p-6 shadow-lg">
        <header className="mb-6 text-center">
          <h1 className="text-xl font-bold text-[hsl(var(--brand))]">Fogistration</h1>
          <span className="text-lg font-semibold">Welcome back</span>
        </header>

        <form className="space-y-4">
          <div className="space-y-2">
            <label htmlFor="email" className="block text-sm font-medium">
              Email
            </label>
            <input
              id="email"
              type="email"
              placeholder="you@example.com"
              className="w-full rounded-lg border border-[hsl(var(--border))] bg-[hsl(var(--input))] px-3 py-2 outline-none ring-[hsl(var(--brand))]/30 transition focus:ring-2"
              value={username}
              onChange={(e) => { setUsername(e.target.value) }}
            />
          </div>

          <div className="space-y-2">
            <div className="flex items-center justify-between">
              <label htmlFor="password" className="text-sm font-medium">
                Password
              </label>
            </div>
            <input
              id="password"
              type="password"
              placeholder="••••••••"
              className="w-full rounded-lg border border-[hsl(var(--border))] bg-[hsl(var(--input))] px-3 py-2 outline-none ring-[hsl(var(--brand))]/30 transition focus:ring-2"
              value={password}
              onChange={(e) => { setPassword(e.target.value) }}
            />
          </div>

          <button
            type="button"
            className="mt-2 w-full rounded-lg bg-[hsl(var(--brand))] px-4 py-2 font-medium text-[hsl(var(--on-brand))] shadow-sm transition hover:opacity-90"
            onClick={signIn}
          >
            Sign in
          </button>

          <div className="relative py-2 text-center text-xs text-[hsl(var(--muted))]">
            <span className="bg-[hsl(var(--panel))] px-2 relative z-10">or continue with</span>
            <span className="absolute left-0 right-0 top-1/2 -z-0 h-px -translate-y-1/2 bg-[hsl(var(--border))]" />
          </div>

          <div className="grid grid-cols-1 gap-3">
            {import.meta.env.DEV &&
              <button
                type="button"
                className="rounded-lg border border-[hsl(var(--border))] bg-[hsl(var(--input))] px-4 py-2 text-sm transition hover:bg-[hsl(var(--hover))]"
                onClick={devLogin}
              >
                Dev Login
              </button>
            }
            <button
              type="button"
              className="rounded-lg border border-[hsl(var(--border))] bg-[hsl(var(--input))] px-4 py-2 text-sm transition hover:bg-[hsl(var(--hover))]"
            >
              OIDC
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
