import { useEffect, useRef } from "react"

export function useDraggable(
  el: HTMLDivElement | null,
  onDelta: (dx: number, dy: number) => void,
  scale = 1
) {
  const cbRef = useRef(onDelta)
  useEffect(() => { cbRef.current = onDelta }, [onDelta])

  const scaleRef = useRef(scale)
  useEffect(() => { scaleRef.current = scale }, [scale])

  useEffect(() => {
    if (!el) return
    let dragging = false
    let sx = 0, sy = 0

    const down = (ev: MouseEvent) => {
      dragging = true
      el.style.cursor = "grabbing"
      sx = ev.clientX; sy = ev.clientY
      ev.preventDefault()
    }
    const move = (ev: MouseEvent) => {
      if (!dragging) return
      const s = scaleRef.current || 1
      cbRef.current((ev.clientX - sx) / s, (ev.clientY - sy) / s)
      sx = ev.clientX; sy = ev.clientY
    }
    const up = () => { dragging = false; el.style.cursor = "grab" }

    el.addEventListener("mousedown", down)
    window.addEventListener("mousemove", move)
    window.addEventListener("mouseup", up)
    return () => {
      el.removeEventListener("mousedown", down)
      window.removeEventListener("mousemove", move)
      window.removeEventListener("mouseup", up)
    }
  }, [el])
}
