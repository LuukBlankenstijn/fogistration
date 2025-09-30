import type { ColumnDef } from "@tanstack/react-table"
import type { HeaderAction, RowAction } from "./types"

export const ACTION_COLUMN_ID = "_actions"

export function createActionsColumn<T>(opts: {
  headerActions?: HeaderAction<T>[]
  rowActions?: RowAction<T>[]
  align?: "left" | "right" | "center"
}): ColumnDef<T> {
  const { headerActions = [], rowActions = [], align = "right" } = opts

  return {
    id: ACTION_COLUMN_ID,
    header: ({ table }) => {
      const selected = table.getSelectedRowModel().rows.map(r => r.original)
      const ctx = { rows: selected, table }

      if (headerActions.length === 1) {
        const a = headerActions[0]
        const dis = a.disabled?.(ctx) ?? selected.length === 0
        return (
          <div className={alignClass(align)}>
            <button
              type="button"
              disabled={dis}
              onClick={() => void a.onClick(ctx)}
              className="rounded px-2 py-1 text-xs disabled:opacity-50"
              style={{ backgroundColor: "hsl(var(--hover))", color: "hsl(var(--fg))" }}
            >
              {a.label}
            </button>
          </div>
        )
      }

      return (
        <div className={alignClass(align)}>
          <details className="relative inline-block">
            <summary
              className="list-none rounded px-2 py-1 text-xs cursor-pointer select-none"
              style={{ backgroundColor: "hsl(var(--hover))", color: "hsl(var(--fg))" }}
              onClick={(e) => {
                if (e.currentTarget instanceof HTMLElement) {
                  e.preventDefault()
                }
              }}
              onMouseDown={(e) => {
                const d = e.currentTarget.parentElement as HTMLDetailsElement
                d.open = !d.open
              }}
            >
              Actions ▾
            </summary>

            <div
              className="absolute right-0 mt-1 w-44 overflow-hidden rounded-md border shadow-lg z-10"
              style={{ backgroundColor: "hsl(var(--panel))", borderColor: "hsl(var(--border))" }}
            >
              <ul className="py-1">
                {headerActions.map((a) => {
                  const dis = a.disabled?.(ctx) ?? selected.length === 0
                  return (
                    <li key={a.label}>
                      <button
                        type="button"
                        disabled={dis}
                        onClick={(e) => {
                          void a.onClick(ctx)
                          const details = e.currentTarget.closest("details") ?? undefined
                          if (details) details.open = false
                        }}
                        className="w-full px-3 py-1.5 text-left text-xs disabled:opacity-50 cursor-pointer"
                        style={{ color: dis ? "hsl(var(--muted))" : "hsl(var(--fg))" }}
                      >
                        {a.label}
                      </button>
                    </li>
                  )
                })}
              </ul>
            </div>
          </details>
        </div >
      )
    },
    cell: ({ row, table }) => {
      const data = row.original
      return (
        <div className={alignClass(align)}>
          {rowActions.map(a => {
            const dis = a.disabled?.({ row: data, table }) ?? false
            return (
              <button
                key={a.label}
                disabled={dis}
                onClick={() => void a.onClick({ row: data, table })}
                className="rounded px-2 py-1 text-xs disabled:opacity-50"
                style={{ backgroundColor: "hsl(var(--hover))", color: "hsl(var(--fg))" }}
                type="button"
              >
                {a.label}
              </button>
            )
          })}
        </div>
      )
    },
    enableSorting: false,
    enableHiding: false,
    size: 1, // let it shrink
    meta: { isActions: true },
  }
}

function alignClass(a: "left" | "right" | "center") {
  return a === "right" ? "text-right" : a === "center" ? "text-center" : "text-left"
}
