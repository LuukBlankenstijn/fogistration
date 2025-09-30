import { Align, type WallpaperLayout } from "@/clients/generated-client";

export const DEFAULT_WALLPAPER_LAYOUT: WallpaperLayout = {
  fontStack: "Noto Sans",
  h: 1080,
  ip: {
    align: Align.LEFT,
    color: "",
    size: 44,
    weight: 600,
    x: 100,
    y: 100,
    display: true,
  },
  teamname: {
    align: Align.CENTER,
    color: "",
    size: 88,
    weight: 800,
    x: 960,
    y: 540,
    display: true,
  },
  w: 1920
}
