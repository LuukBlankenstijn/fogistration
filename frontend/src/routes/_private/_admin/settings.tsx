import Settings from '@/components/Settings'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_private/_admin/settings')({
  component: Settings,
})

