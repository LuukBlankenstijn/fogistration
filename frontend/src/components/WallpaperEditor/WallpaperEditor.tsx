import { useEffect, useState } from "react"
import type { Layout, WallpaperEditorProps } from "./types/layout"
import Settings from "./Settings"
import Preview from "./Preview"
import { useGetWallpaperConfigQuery, useGetWallpaperQuery, useWallpaperConfigMutation, useWallpaperMutation } from "@/query/wallpaper"

export default function WallpaperEditor({ contestId }: WallpaperEditorProps) {
  const { data: initialLayout } = useGetWallpaperConfigQuery(contestId)
  const { data: wallpaperBlob } = useGetWallpaperQuery(contestId)
  const { mutate: saveConfig, isPending: configIsPending } = useWallpaperConfigMutation()
  const { mutate: saveWallpaper, isPending: wallpaperIsPending } = useWallpaperMutation(contestId)

  useEffect(() => {
    onBGFile(wallpaperBlob ?? undefined)
  }, [])


  const onBGFile = (file?: File | Blob) => {
    if (!file) {
      if (bgUrl) URL.revokeObjectURL(bgUrl)
      setBgUrl(undefined)
      return
    }
    const url = URL.createObjectURL(file)
    if (bgUrl) URL.revokeObjectURL(bgUrl)
    setBgUrl(url)
  }


  const [layout, setLayout] = useState<Layout>(initialLayout)
  const [bgUrl, setBgUrl] = useState<string | undefined>(undefined)

  const save = () => {
    saveConfig({ layout, contestId })
    saveWallpaper({ url: bgUrl })
  }
  return (
    <main className="flex h-[100svh] gap-4 p-4 text-[hsl(var(--fg))] overflow-hidden">
      <Preview layout={layout} setLayout={setLayout} bgUrl={bgUrl} />

      <Settings
        layout={layout}
        setLayout={setLayout}
        bgUrl={bgUrl}
        setBg={onBGFile}
        save={save}
        isPending={wallpaperIsPending && configIsPending}
      />
    </main>
  )
}
