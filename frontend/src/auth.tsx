import { createContext, use, useCallback, useMemo } from "react"
import { getCurrentUserOptions } from "./clients/generated-client/@tanstack/react-query.gen"
import { useQuery } from "@tanstack/react-query"
import { useDevLogin, useLogin, useLogout } from "./query/auth"
import type { User } from "./clients/generated-client"

export interface AuthState {
  isAuthenticated: boolean
  user: User | null
  login: (username: string, password: string) => void
  logout: () => void
  devLogin: () => void
}

const AuthContext = createContext<AuthState | undefined>(undefined);

export const AuthProvider = ({ children }: { children: React.ReactNode }) => {
  const userOptions = useMemo(() => getCurrentUserOptions(), [])
  userOptions.retry = false
  const { data: currentUserResponse, isLoading, isError } = useQuery(userOptions)
  const { user, isAuthenticated } = useMemo(() => ({
    user: currentUserResponse,
    isAuthenticated: currentUserResponse ?? false
  }), [currentUserResponse])
  const { mutate: mutateLogin } = useLogin()
  const { mutate: mutateLogout } = useLogout()
  const { mutate: mutateDevLogin } = useDevLogin()
  const login = useCallback((username: string, password: string) => {
    mutateLogin({
      body: {
        username,
        password
      }
    })
  }, [mutateLogin])

  const devLogin = useCallback(() => {
    mutateDevLogin({})
  }, [mutateDevLogin])

  const logout = useCallback(() => {
    mutateLogout({})
  }, [])


  const value = useMemo((): AuthState => ({
    isAuthenticated: isAuthenticated && !isError,
    user: user ?? null,
    login,
    logout,
    devLogin,
  }), [currentUserResponse, isError, login, logout])

  if (isLoading) {
    return (
      <div>
        Loading...
      </div>
    )
  }

  return (
    <AuthContext value={value}>
      {children}
    </AuthContext>
  )
}

export const useAuth = () => {
  const context = use(AuthContext)
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider")
  }
  return context
}
