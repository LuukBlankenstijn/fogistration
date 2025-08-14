interface DropdownProps<T> {
  label?: string
  value: string | undefined
  onChange: (value: string | null) => void
  options: T[]
  valueGenerator: (row: T) => string
  labelGenerator: (row: T) => string
  showGenerator: (row: T, value: string | undefined) => boolean,
  placeholder?: string
  required?: boolean
}

export function Dropdown<T>({
  label,
  value,
  onChange,
  options,
  valueGenerator,
  labelGenerator,
  showGenerator,
  placeholder = "Select…",
  required = false,
}: DropdownProps<T>) {
  return (
    <div className="flex flex-col gap-2">
      {label && <label className="text-sm font-medium">{label}</label>}
      <select
        value={value ?? ""}
        onChange={(e) => { onChange(e.target.value === "" ? null : e.target.value); }}
        className="rounded-lg border border-[hsl(var(--border))] bg-[hsl(var(--input))] px-3 py-2 text-sm text-[hsl(var(--fg))] shadow-sm focus:border-[hsl(var(--brand))] focus:outline-none focus:ring-2 focus:ring-[hsl(var(--brand))]"
      >
        <option value="" disabled={required} className="text-[hsl(var(--muted))]">
          {placeholder}
        </option>
        {options.map((opt) => {
          return (showGenerator(opt, value) &&
            <option
              key={valueGenerator(opt)}
              value={valueGenerator(opt)}
            >
              {labelGenerator(opt)}
            </option>
          )
        })}
      </select>
    </div>
  )
}

