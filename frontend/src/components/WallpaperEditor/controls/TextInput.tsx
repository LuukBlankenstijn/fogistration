export function TextInput({ value, onChange }: { value: string; onChange: (v: string) => void }) {
  return (
    <input
      type="text"
      value={value}
      onChange={(e) => { onChange(e.target.value); }}
      className="rounded-lg border border-[hsl(var(--border))] bg-[hsl(var(--input))] px-3 py-2"
    />
  )
}
