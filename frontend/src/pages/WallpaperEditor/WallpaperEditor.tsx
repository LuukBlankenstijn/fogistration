import { useState } from "react"
import Settings from "./Settings"
import Preview from "./Preview"
import { useGetWallpaperConfigQuery, useGetWallpaperQuery as useGetWallpaperFileQuery, useWallpaperMutation } from "@/query/wallpaper"
import type { WallpaperLayout } from "@/clients/generated-client"

interface WallpaperEditorProps {
  id: number
}

export default function WallpaperEditor({ id }: WallpaperEditorProps) {
  const { data: wallpaperLayout } = useGetWallpaperConfigQuery(id)
  const { data: wallpaperBlob } = useGetWallpaperFileQuery(id)
  const { mutate, isPending } = useWallpaperMutation(id)

  const [layout, setLayout] = useState<WallpaperLayout>(wallpaperLayout)
  const [file, setFile] = useState<Blob | File | null>(wallpaperBlob)

  const save = () => {
    mutate({ layout, file })
  }
  return (
    <main className="flex h-[100svh] gap-4 p-4 text-[hsl(var(--fg))] overflow-hidden">
      <Preview layout={layout} setLayout={setLayout} file={file} />

      <Settings
        layout={layout}
        setLayout={setLayout}
        file={file}
        setFile={(file) => { setFile(file ?? null) }}
        save={save}
        isPending={isPending}
      />
    </main>
  )
}
