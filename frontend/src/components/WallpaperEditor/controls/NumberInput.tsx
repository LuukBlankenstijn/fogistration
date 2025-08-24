export function NumberInput({ value, onChange, step = 10, className }: { value: number; onChange: (v: number) => void; step?: number, className?: string }) {
  return (
    <input
      type="number"
      value={value}
      step={step}
      onChange={(e) => { onChange(Number.isFinite(+e.target.value) ? +e.target.value : value); }}
      className={(className ?? "") + " w-24 min-w-0 rounded-lg border border-[hsl(var(--border))] bg-[hsl(var(--input))] px-3 py-2"}
    />
  )
}
