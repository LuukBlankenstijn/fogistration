import { ChevronLeft } from "lucide-react"
import { useNavigate } from "@tanstack/react-router"
import { useEffect } from "react"
import type { Route } from "@/routeTree"


export function GoBackButton({ to }: { to: Route }) {
  const navigate = useNavigate()

  const goBack = () => {
    void navigate({ to: to })
  }

  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") {
        goBack()
      }
    }
    window.addEventListener("keydown", onKey)
    return () => { window.removeEventListener("keydown", onKey); }
  }, [])

  return (
    <button
      type="button"
      onClick={goBack}
      className="rounded-full p-2 hover:bg-[hsl(var(--hover))] transition-colors cursor-pointer"
      aria-label="Go back"
    >
      <ChevronLeft className="h-5 w-5 text-[hsl(var(--fg))]" />
    </button>
  )
}
