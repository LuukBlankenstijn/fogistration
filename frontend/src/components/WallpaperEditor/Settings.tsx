import { useRef } from "react";
import { Card } from "../Card";
import { NumberInput } from "./controls/NumberInput";
import { TextInput } from "./controls/TextInput";
import { Align, type WallpaperLayout } from "@/clients/generated-client";

interface SetttingsProps {
  layout: WallpaperLayout,
  file: File | Blob | null
  setLayout: (value: React.SetStateAction<WallpaperLayout>) => void
  setFile: (file?: File | Blob) => void
  save: () => void
  isPending: boolean
}

const Settings = ({ layout, setLayout, setFile, save, isPending }: SetttingsProps) => {
  const fileRef = useRef<HTMLInputElement>(null)

  const onBGFile = (file?: File) => {
    if (!file) {
      if (fileRef.current) fileRef.current.value = ""
    }
    setFile(file)
  }

  return (
    <Card className="w-70 h-full shrink-0 overflow-auto">
      <h2 className="text-lg font-semibold">Wallpaper editor</h2>
      <div className="mt-3 space-y-4">

        <div>
          <label className="text-sm text-[hsl(var(--muted))]">Canvas size</label>
          <div className="mt-1 flex items-center gap-2">
            <NumberInput value={layout.w} onChange={(v) => { setLayout(l => ({ ...l, w: v })); }} />
            <span>×</span>
            <NumberInput value={layout.h} onChange={(v) => { setLayout(l => ({ ...l, h: v })); }} />
          </div>
        </div>

        <div>
          <label className="text-sm text-[hsl(var(--muted))]">Background PNG</label>
          <div className="mt-1 flex items-center gap-2 min-w-0">
            {/* Hidden file input */}
            <input
              ref={fileRef}
              type="file"
              accept="image/png"
              className="hidden"
              onChange={(e) => { onBGFile(e.target.files?.[0]); }}
            />

            {/* Browse button */}
            <button
              type="button"
              onClick={() => fileRef.current?.click()}
              className="shrink-0 rounded-lg border border-[hsl(var(--border))] bg-[hsl(var(--input))] px-3 py-2 text-sm hover:bg-[hsl(var(--hover))]"
            >
              Browse
            </button>

            {/* Clear button */}
            <button
              type="button"
              className="shrink-0 rounded-lg border border-[hsl(var(--border))] bg-[hsl(var(--input))] px-3 py-2 text-sm hover:bg-[hsl(var(--hover))]"
              onClick={() => { onBGFile(undefined); }}
            >
              Clear
            </button>
          </div>
        </div>


        <div>
          <label className="text-sm text-[hsl(var(--muted))]">Font stack</label>
          <div className="mt-1 flex items-center gap-2">
            <TextInput value={layout.fontStack} onChange={(v) => { setLayout(l => ({ ...l, fontStack: v })); }} />
          </div>
        </div>

        <fieldset className="grid gap-2 rounded-xl border border-[hsl(var(--border))] p-3">
          <legend className="px-1 text-sm">teamname</legend>

          <div className="grid grid-cols-2 items-center">
            <label className="text-sm text-[hsl(var(--muted))]">Size</label>
            <NumberInput
              value={layout.teamname.size}
              onChange={(v) => { setLayout(l => ({ ...l, teamname: { ...l.teamname, size: v } })); }}
              className="justify-self-end"
            />
          </div>

          <div className="grid grid-cols-2 items-center">
            <label className="text-sm text-[hsl(var(--muted))]">Weight</label>
            <NumberInput
              value={layout.teamname.weight}
              step={100}
              onChange={(v) => { setLayout(l => ({ ...l, teamname: { ...l.teamname, weight: v } })); }}
              className="justify-self-end"
            />
          </div>

          <div className="grid grid-cols-2 items-center">
            <label className="text-sm text-[hsl(var(--muted))]">Color</label>
            <input
              type="color"
              className="justify-self-end h-9 w-12 rounded-md border border-[hsl(var(--border))] bg-[hsl(var(--input))]"
              value={layout.teamname.color}
              onChange={(e) => { setLayout(l => ({ ...l, teamname: { ...l.teamname, color: e.target.value } })); }}
            />
          </div>

          <div className="grid grid-cols-2 items-center">
            <label className="text-sm text-[hsl(var(--muted))]">Align</label>
            <select
              className="justify-self-end rounded-lg border border-[hsl(var(--border))] bg-[hsl(var(--input))] px-2 py-2 text-sm"
              value={layout.teamname.align}
              onChange={(e) => { setLayout(l => ({ ...l, teamname: { ...l.teamname, align: e.target.value as Align } })); }}
            >
              <option value={Align.LEFT}>Left</option>
              <option value={Align.CENTER}>Center</option>
              <option value={Align.RIGHT}>Right</option>
            </select>
          </div>
        </fieldset>

        <fieldset className="grid gap-2 rounded-xl border border-[hsl(var(--border))] p-3">
          <legend className="px-1 text-sm">ip</legend>

          <div className="grid grid-cols-2 items-center">
            <label className="text-sm text-[hsl(var(--muted))]">Size</label>
            <NumberInput
              value={layout.ip.size}
              onChange={(v) => { setLayout(l => ({ ...l, ip: { ...l.ip, size: v } })); }}
              className="justify-self-end"
            />
          </div>

          <div className="grid grid-cols-2 items-center">
            <label className="text-sm text-[hsl(var(--muted))]">Weight</label>
            <NumberInput
              value={layout.ip.weight}
              step={100}
              onChange={(v) => { setLayout(l => ({ ...l, ip: { ...l.ip, weight: v } })); }}
              className="justify-self-end"
            />
          </div>

          <div className="grid grid-cols-2 items-center">
            <label className="text-sm text-[hsl(var(--muted))]">Color</label>
            <input
              type="color"
              className="justify-self-end h-9 w-12 rounded-md border border-[hsl(var(--border))] bg-[hsl(var(--input))]"
              value={layout.ip.color}
              onChange={(e) => { setLayout(l => ({ ...l, ip: { ...l.ip, color: e.target.value } })); }}
            />
          </div>

          <div className="grid grid-cols-2 items-center">
            <label className="text-sm text-[hsl(var(--muted))]">Align</label>
            <select
              className="justify-self-end rounded-lg border border-[hsl(var(--border))] bg-[hsl(var(--input))] px-2 py-2 text-sm"
              value={layout.ip.align}
              onChange={(e) => { setLayout(l => ({ ...l, ip: { ...l.ip, align: e.target.value as Align } })); }}
            >
              <option value={Align.LEFT}>Left</option>
              <option value={Align.CENTER}>Center</option>
              <option value={Align.RIGHT}>Right</option>
            </select>
          </div>
        </fieldset>


        <button
          type="button"
          onClick={() => { save(); }}
          className="w-full shrink-0 rounded-lg border border-[hsl(var(--border))] bg-[hsl(var(--input))] px-3 py-2 text-sm hover:bg-[hsl(var(--hover))]"
        >
          {isPending ? "Saving..." : "Save"}
        </button>

      </div>
    </Card>
  )
}

export default Settings
