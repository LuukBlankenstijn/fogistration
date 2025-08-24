import { useCallback, useRef } from "react";
import { Card } from "../Card";
import { useElementSize } from "./hooks/useElementSize";
import type { Layout } from "./types/layout";
import { useDraggable } from "./hooks/useDraggable";
import { GoBackButton } from "../GoBack";

interface PreviewProps {
  layout: Layout,
  setLayout: (value: React.SetStateAction<Layout>) => void
  bgUrl?: string
}

export default function Preview({ layout, setLayout, bgUrl }: PreviewProps) {
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

  const teamRef = useRef<HTMLDivElement>(null)
  const ipRef = useRef<HTMLDivElement>(null)


  useDraggable(teamRef, onTeamDrag, scale)
  useDraggable(ipRef, onIpDrag, scale)

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
          {bgUrl ? (
            <img src={bgUrl} alt="bg" className="absolute inset-0 h-full w-full object-cover" />
          ) : (
            <div className="absolute inset-0 grid place-items-center text-[hsl(var(--muted))]">
              <div className="text-center text-sm">No background — choose a file</div>
            </div>
          )}

          <div
            ref={teamRef}
            className="absolute cursor-grab whitespace-pre drop-shadow-[0_2px_6px_rgba(0,0,0,0.60)]"
            style={labelStyle(layout, "teamname")}
          >
            {"{{ teamname }}"}
          </div>

          <div
            ref={ipRef}
            className="absolute cursor-grab whitespace-pre drop-shadow-[0_2px_6px_rgba(0,0,0,0.60)]"
            style={labelStyle(layout, "ip")}
          >
            {"{{ ip }}"}
          </div>
        </div>
      </div>
    </Card>
  )
}


function labelStyle(l: Layout, key: "teamname" | "ip"): React.CSSProperties {
  const spec = l[key]
  const tx = spec.align === "center" ? "-50%" : spec.align === "right" ? "-100%" : "0"
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
