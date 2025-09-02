export function BooleanInput({
  value,
  onChange,
  className,
}: {
  value: boolean
  onChange: (v: boolean) => void
  className?: string
}) {
  return (
    <button
      type="button"
      role="switch"
      aria-checked={value}
      onClick={() => { onChange(!value); }}
      className={
        (className ?? "") +
        " relative inline-flex h-6 w-12 items-center rounded-full border border-[hsl(var(--border))] transition-colors " +
        (value
          ? "bg-[hsl(var(--input))]"
          : "bg-[hsl(var(--input-disabled))]")
      }
    >
      <span
        className={
          "inline-block h-4 w-4 transform rounded-full bg-white shadow transition-transform " +
          (value ? "translate-x-6" : "translate-x-1")
        }
      />
    </button>
  )
}
