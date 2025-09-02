export function NumberInput({
  value,
  onChange,
  step = 10,
  className,
  disabled = false,
}: {
  value: number
  onChange: (v: number) => void
  step?: number
  className?: string
  disabled?: boolean
}) {
  return (
    <input
      type="number"
      value={value}
      step={step}
      disabled={disabled}
      onChange={(e) => {
        if (!disabled) {
          const next = +e.target.value
          onChange(Number.isFinite(next) ? next : value)
        }
      }}
      className={
        (className ?? "") +
        " w-24 min-w-0 rounded-lg border px-3 py-2 " +
        (disabled
          ? "cursor-not-allowed border-[hsl(var(--border))] bg-[hsl(var(--input-disabled))] text-[hsl(var(--fg-disabled))]"
          : "border-[hsl(var(--border))] bg-[hsl(var(--input))] text-[hsl(var(--fg))]")
      }
    />
  )
}
