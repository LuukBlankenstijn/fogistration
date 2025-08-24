import WallpaperEditor from '@/components/WallpaperEditor/WallpaperEditor'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_private/wallpaper/$contestId')({
  component: RouteComponent,
})

function RouteComponent() {
  const { contestId } = Route.useParams()

  return <WallpaperEditor contestId={contestId} />
}
