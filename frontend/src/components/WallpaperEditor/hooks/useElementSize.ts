import { useEffect, useRef, useState } from "react"

export const useElementSize = <T extends HTMLElement>() => {
  const ref = useRef<T | null>(null)
  const [size, setSize] = useState({ w: 0, h: 0 })
  useEffect(() => {
    const el = ref.current
    if (!el) return
    const ro = new ResizeObserver(([entry]) => {
      const cr = entry.contentRect
      setSize({ w: cr.width, h: cr.height })
    })
    ro.observe(el)
    return () => { ro.disconnect(); }
  }, [])
  return { ref, size }
}
