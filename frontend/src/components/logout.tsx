import { useAuth } from "@/auth"

type LogoutButtonProps = React.ButtonHTMLAttributes<HTMLButtonElement>

export function LogoutButton({ className = '', children = 'Logout', ...rest }: LogoutButtonProps) {
  const { logout } = useAuth()
  return (
    <button
      type="button"
      {...rest}
      className={`rounded-lg border border-[hsl(var(--border))] bg-[hsl(var(--input))] px-3 py-2 text-sm transition hover:bg-[hsl(var(--hover))] ${className}`}
      onClick={logout}
    >
      {children}
    </button>
  )
}
