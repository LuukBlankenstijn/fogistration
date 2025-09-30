import { flexRender, type RowData, type Table } from "@tanstack/react-table"

interface StyledTableProps<T> {
  tableData: Table<T>
  title?: string
  loading?: boolean
}

const StyledTable = <T extends RowData>({ tableData, title, loading = false }: StyledTableProps<T>) => {
  return (
    <div>
      <header className="flex items-baseline justify-between">
        <h2 className="text-lg font-semibold" style={{ color: 'hsl(var(--fg))' }}>{title}</h2>
        <span className="text-xs" style={{ color: 'hsl(var(--muted))' }}>{tableData.getRowCount()} total</span>
      </header>

      <div className="overflow-x-auto rounded-lg border"
        style={{ backgroundColor: "hsl(var(--panel))", borderColor: "hsl(var(--border))" }}>
        <table className="w-full border-collapse text-sm">
          <thead style={{ backgroundColor: "hsl(var(--hover))" }}>
            {tableData.getHeaderGroups().map(hg => (
              <tr key={hg.id}>
                {hg.headers.map(header => (
                  <th key={header.id}
                    onClick={header.column.getToggleSortingHandler()}
                    className="px-4 py-2 text-left font-medium select-none cursor-pointer"
                    style={{ color: "hsl(var(--fg))", borderBottom: "1px solid hsl(var(--border))" }}>
                    {header.isPlaceholder ? null : flexRender(header.column.columnDef.header, header.getContext())}
                    {header.column.getIsSorted() === 'asc' && ' ▲'}
                    {header.column.getIsSorted() === 'desc' && ' ▼'}
                  </th>
                ))}
              </tr>
            ))}
          </thead>

          <tbody>
            {tableData.getRowModel().rows.map(row => (
              <tr key={row.id} className="transition-colors" style={{ color: "hsl(var(--fg))" }}>
                {row.getVisibleCells().map(cell => (
                  <td key={cell.id} className="px-4 py-2" style={{ borderBottom: "1px solid hsl(var(--border))" }}>
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </td>
                ))}
              </tr>
            ))}
            {(tableData.getRowModel().rows.length === 0 || loading) && (
              <tr>
                <td colSpan={tableData.getAllColumns().length} className="px-4 py-6 text-center" style={{ color: "hsl(var(--muted))" }}>
                  {loading ? "Loading..." : "No clients found"}
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
