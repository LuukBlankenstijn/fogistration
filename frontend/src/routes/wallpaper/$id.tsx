import WallpaperEditor from '@/components/WallpaperEditor/WallpaperEditor'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/wallpaper/$id')({
  component: RouteComponent,
  params: {
    parse: (raw) => {
      const id = Number(raw.id)
      if (!Number.isInteger(id)) {
        throw new Error(`Invalid id: ${raw.id}`)
      }
      return { id }
    },
    stringify: (parsed) => ({ id: String(parsed.id) })
  }
})

function RouteComponent() {
  const { id } = Route.useParams()

  return <WallpaperEditor id={id} />
}
