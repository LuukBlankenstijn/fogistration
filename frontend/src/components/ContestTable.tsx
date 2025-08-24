import { getAllContestsOptions } from "@/clients/generated-client/@tanstack/react-query.gen";
import { useSuspenseQuery } from "@tanstack/react-query";
import { type SortingState, type ColumnDef, useReactTable, getCoreRowModel, getSortedRowModel, flexRender } from "@tanstack/react-table";
import { useState, useMemo } from "react";
import type { ModelsContest } from "@/clients/generated-client";
import { Link } from "@tanstack/react-router";

export function ContestTable() {
  const { data: contests } = useSuspenseQuery(getAllContestsOptions())

  const [sorting, setSorting] = useState<SortingState>([])

  const columns = useMemo<ColumnDef<ModelsContest>[]>(() => [
    { accessorKey: "id", header: "ID" },
    { accessorKey: "name", header: "Name" },
    {
      accessorKey: "startTime",
      header: "Start time",
      cell: ({ getValue }) => {
        const d = new Date(getValue() as string)
        return d.toLocaleString("nl-NL", {
          year: "numeric",
          month: "2-digit",
          day: "2-digit",
          hour: "2-digit",
          minute: "2-digit"
        })
      }
    },
    {
      accessorKey: "endTime",
      header: "End time",
      cell: ({ getValue }) => {
        const d = new Date(getValue() as string)
        return d.toLocaleString("nl-NL", {
          year: "numeric",
          month: "2-digit",
          day: "2-digit",
          hour: "2-digit",
          minute: "2-digit"
        })
      }
    },
    {
      id: "actions",
      header: "wallpaper",
      enableSorting: false,
      cell: ({ row }) => {
        return (
          <Link to="/wallpaper/$contestId" params={{ contestId: row.original.id }}>
            <button
              type="button"
              className="inline-flex items-center rounded-md border px-3 py-1.5 text-xs font-medium transition focus:outline-none focus:ring-1 focus:ring-offset-1"
              style={{
                backgroundColor: "hsl(var(--input))",
                color: "hsl(var(--fg))",
                borderColor: "hsl(var(--border))",
              }}
              onMouseEnter={(e) => (e.currentTarget.style.backgroundColor = "hsl(var(--hover))")}
              onMouseLeave={(e) => (e.currentTarget.style.backgroundColor = "hsl(var(--input))")}
            >
              Edit
            </button>
          </Link>
        )
      },
      size: 1, // keeps it compact (optional, if you use column sizing)
    },
  ], [])

  const table = useReactTable({
    data: contests,
    columns,
    state: { sorting },
    onSortingChange: setSorting,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
  })

  return (
    <div>
      <header className="flex items-baseline justify-between">
        <h2 className="text-lg font-semibold" style={{ color: 'hsl(var(--fg))' }}>Contests</h2>
        <span className="text-xs" style={{ color: 'hsl(var(--muted))' }}>{contests.length} total</span>
      </header>

      <div className="overflow-x-auto rounded-lg border"
        style={{ backgroundColor: "hsl(var(--panel))", borderColor: "hsl(var(--border))" }}>
        <table className="w-full border-collapse text-sm">
          <thead style={{ backgroundColor: "hsl(var(--hover))" }}>
            {table.getHeaderGroups().map(hg => (
              <tr key={hg.id}>
                {hg.headers.map(header => (
                  <th key={header.id}
                    onClick={header.column.getToggleSortingHandler()}
                    className="px-4 py-2 text-left font-medium select-none cursor-pointer"
                    style={{ color: "hsl(var(--fg))", borderBottom: "1px solid hsl(var(--border))" }}>
                    {header.isPlaceholder ? null : flexRender(header.column.columnDef.header, header.getContext())}
                    {/* tiny sort indicator */}
                    {header.column.getIsSorted() === 'asc' && ' ▲'}
                    {header.column.getIsSorted() === 'desc' && ' ▼'}
                  </th>
                ))}
              </tr>
            ))}
          </thead>

          <tbody>
            {table.getRowModel().rows.map(row => (
              <tr key={row.id} className="transition-colors" style={{ color: "hsl(var(--fg))" }}>
                {row.getVisibleCells().map(cell => (
                  <td key={cell.id} className="px-4 py-2" style={{ borderBottom: "1px solid hsl(var(--border))" }}>
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </td>
                ))}
              </tr>
            ))}
            {table.getRowModel().rows.length === 0 && (
              <tr>
                <td colSpan={columns.length} className="px-4 py-6 text-center" style={{ color: "hsl(var(--muted))" }}>
                  No clients found.
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  )
}
