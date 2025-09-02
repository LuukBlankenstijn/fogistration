import { useCallback, useMemo, useState } from "react";
import { Card } from "../Card";
import { useElementSize } from "./hooks/useElementSize";
import { useDraggable } from "./hooks/useDraggable";
import { GoBackButton } from "../GoBack";
import { Align, type WallpaperLayout } from "@/clients/generated-client";

interface PreviewProps {
  layout: WallpaperLayout,
  setLayout: (value: React.SetStateAction<WallpaperLayout>) => void
  file: File | Blob | null
}

export default function Preview({ layout, setLayout, file }: PreviewProps) {
  const { ref: previewRef, size: previewSize } = useElementSize<HTMLDivElement>()

  const onTeamDrag = useCallback((dx: number, dy: number) => {
    setLayout(l => ({
      ...l,
      teamname: {
        ...l.teamname,
        x: Math.round(l.teamname.x + dx),
        y: Math.round(l.teamname.y + dy),
      },
    }))
  }, [])

  const url = useMemo(() => {
    return file ? URL.createObjectURL(file) : undefined
  }, [file])

  const onIpDrag = useCallback((dx: number, dy: number) => {
    setLayout(l => ({
      ...l,
      ip: {
        ...l.ip,
        x: Math.round(l.ip.x + dx),
        y: Math.round(l.ip.y + dy),
      },
    }))
  }, [])

  const scale = Math.min(
    (previewSize.w / layout.w) || 1,
    (previewSize.h / layout.h) || 1,
    1
  )

  const [teamEl, setTeamEl] = useState<HTMLDivElement | null>(null)
  const [ipEl, setIpEl] = useState<HTMLDivElement | null>(null)

  useDraggable(teamEl, onTeamDrag, scale)
  useDraggable(ipEl, onIpDrag, scale)

  return (
    <Card className="min-w-0 flex-1 h-full overflow-hidden p-4">
      <div className="mb-2 flex items-center gap-2">
        <GoBackButton />
        <h2 className="text-lg font-semibold">Wallpaper preview</h2>
      </div>
      {/* measuring target */}
      <div ref={previewRef} className="h-full w-full grid place-items-center overflow-hidden">
        {/* scaled stage */}
        <div
          className="relative select-none rounded-2xl border border-[hsl(var(--border))] bg-black shadow-sm"
          style={{
            width: layout.w,
            height: layout.h,
            transform: `scale(${scale.toString()})`,
            transformOrigin: "top left",
          }}
        >
          {file ? (
            <img src={url} alt="bg" className="absolute inset-0 h-full w-full object-cover" />
          ) : (
            <div className="absolute inset-0 grid place-items-center text-[hsl(var(--muted))]">
              <div className="text-center text-sm">No background — choose a file</div>
            </div>
          )}

          {layout.teamname.display &&
            <div
              ref={setTeamEl}
              className="absolute cursor-grab whitespace-pre drop-shadow-[0_2px_6px_rgba(0,0,0,0.60)]"
              style={labelStyle(layout, "teamname")}
            >
              {"{{ teamname }}"}
            </div>
          }

          {layout.ip.display &&
            <div
              ref={setIpEl}
              className="absolute cursor-grab whitespace-pre drop-shadow-[0_2px_6px_rgba(0,0,0,0.60)]"
              style={labelStyle(layout, "ip")}
            >
              {"{{ ip }}"}
            </div>
          }
        </div>
      </div>
    </Card>
  )
}


function labelStyle(l: WallpaperLayout, key: "teamname" | "ip"): React.CSSProperties {
  const spec = l[key]
  const tx = spec.align === Align.CENTER ? "-50%" : spec.align === Align.RIGHT ? "-100%" : "0"
  return {
    left: spec.x,
    top: spec.y,
    color: spec.color,
    fontFamily: l.fontStack,
    fontWeight: spec.weight,
    fontSize: spec.size,
    textAlign: spec.align,
    transform: `translate(${tx}, -100%)`,
  }
}
