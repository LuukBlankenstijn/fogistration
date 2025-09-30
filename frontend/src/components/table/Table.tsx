import * as React from "react"
import { flexRender, type RowData, type Table as RTTable } from "@tanstack/react-table"

interface StyledTableProps<T extends RowData> {
  table: RTTable<T>
  title?: string
  loading?: boolean
}

function IndeterminateCheckbox(
  props: React.InputHTMLAttributes<HTMLInputElement> & { indeterminate?: boolean }
) {
  const { indeterminate, ...rest } = props
  const ref = React.useRef<HTMLInputElement>(null)
  React.useEffect(() => {
    if (ref.current) ref.current.indeterminate = !!indeterminate && !rest.checked
  }, [indeterminate, rest.checked])
  return <input ref={ref} type="checkbox" {...rest} />
}

const StyledTable = <T extends RowData>({
  table,
  title,
  loading = false,
}: StyledTableProps<T>) => {
  const rowCount = table.getRowCount()
  const rows = table.getRowModel().rows
  const empty = rows.length === 0 || loading

  // detect if we have an actions column
  const hasActions = table.getAllColumns().some((c) => c.id === "_actions")

  return (
    <div>
      <header className="flex items-baseline justify-between">
        <h2 className="text-lg font-semibold" style={{ color: "hsl(var(--fg))" }}>
          {title}
        </h2>
        <span className="text-xs" style={{ color: "hsl(var(--muted))" }}>
          {rowCount} total
        </span>
      </header>

      <div
        className="overflow-x-auto rounded-lg border"
        style={{ backgroundColor: "hsl(var(--panel))", borderColor: "hsl(var(--border))" }}
      >
        <table className="w-full border-collapse text-sm">
          <thead style={{ backgroundColor: "hsl(var(--hover))" }}>
            {table.getHeaderGroups().map((hg) => (
              <tr key={hg.id}>
                {hasActions && (
                  <th
                    className="px-2 py-2 text-center"
                    style={{ borderBottom: "1px solid hsl(var(--border))" }}
                  >
                    <IndeterminateCheckbox
                      checked={table.getIsAllRowsSelected()}
                      indeterminate={table.getIsSomeRowsSelected()}
                      onChange={table.getToggleAllRowsSelectedHandler()}
                    />
                  </th>
                )}
                {hg.headers.map((header) => (
                  <th
                    key={header.id}
                    onClick={header.column.getToggleSortingHandler()}
                    className="px-4 py-2 text-left font-medium select-none cursor-pointer"
                    style={{ color: "hsl(var(--fg))", borderBottom: "1px solid hsl(var(--border))" }}
                  >
                    {header.isPlaceholder ? null : flexRender(header.column.columnDef.header, header.getContext())}
                    {header.column.getIsSorted() === "asc" && " ▲"}
                    {header.column.getIsSorted() === "desc" && " ▼"}
                  </th>
                ))}
              </tr>
            ))}
          </thead>

          <tbody>
            {rows.map((row) => (
              <tr key={row.id} className="transition-colors" style={{ color: "hsl(var(--fg))" }}>
                {hasActions && (
                  <td className="px-2 py-2 text-center" style={{ borderBottom: "1px solid hsl(var(--border))" }}>
                    <IndeterminateCheckbox
                      checked={row.getIsSelected()}
                      indeterminate={row.getIsSomeSelected()}
                      onChange={row.getToggleSelectedHandler()}
                    />
                  </td>
                )}
                {row.getVisibleCells().map((cell) => (
                  <td
                    key={cell.id}
                    className="px-4 py-2"
                    style={{ borderBottom: "1px solid hsl(var(--border))" }}
                  >
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </td>
                ))}
              </tr>
            ))}

            {empty && (
              <tr>
                <td
                  colSpan={table.getAllColumns().length + (hasActions ? 1 : 0)}
                  className="px-4 py-6 text-center"
                  style={{ color: "hsl(var(--muted))" }}
                >
                  {loading ? "Loading..." : "No rows found"}
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  )
}

export default StyledTable

