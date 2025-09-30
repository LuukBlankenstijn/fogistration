import type { Route } from "@/routeTree"
import { Link, useRouterState } from "@tanstack/react-router"
import { Home, Users, Trophy, type LucideProps } from "lucide-react"
import type { ForwardRefExoticComponent, RefAttributes } from "react"



export type LucideIcon = ForwardRefExoticComponent<
  Omit<LucideProps, "ref"> & RefAttributes<SVGSVGElement>
>
interface SidebarItem {
  to: Route
  label: string
  Icon: LucideIcon
}
const items: SidebarItem[] = [
  { to: "/clients", label: "Clients", Icon: Users },
  { to: "/teams", label: "Teams", Icon: Home },
  { to: "/contests", label: "Contests", Icon: Trophy },
]

export default function Sidebar() {
  const { location } = useRouterState()
  return (
    <aside
      className="hidden md:flex md:w-40 md:flex-col md:gap-2 md:p-3 md:sticky md:top-0 md:h-dvh border-r"
      style={{ backgroundColor: "hsl(var(--panel))", borderColor: "hsl(var(--border))" }}
    >
      <div className="px-2 py-1 text-sm font-semibold select-none" style={{ color: "hsl(var(--fg))" }}>
        Fogistration
      </div>
      <nav className="flex flex-col gap-1">
        {items.map(({ to, label, Icon }) => {
          const active = location.pathname.startsWith(to)
          return (
            <Link
              key={to}
              to={to}
              className={[
                "flex items-center gap-2 rounded-lg px-3 py-2 text-sm transition-colors",
                active ? "font-medium" : "opacity-80 hover:opacity-100",
              ].join(" ")}
              style={{
                color: "hsl(var(--fg))",
                backgroundColor: active ? "hsl(var(--hover))" : "transparent",
              }}
            >
              <Icon size={16} />
              {label}
            </Link>
          )
        })}
      </nav>
    </aside>
  )
}
