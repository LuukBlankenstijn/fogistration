import type { Contest } from "@/clients/generated-client"
import { listContestsOptions } from "@/clients/generated-client/@tanstack/react-query.gen"
import StyledTable from "@/components/table/Table";
import { useSuspenseQuery } from "@tanstack/react-query"
import { Link } from "@tanstack/react-router"
import { type SortingState, type ColumnDef, useReactTable, getCoreRowModel, getSortedRowModel } from "@tanstack/react-table"
import { useState, useMemo } from "react"

const ContestPage = () => {
  const { data: contests } = useSuspenseQuery(listContestsOptions())

  const [sorting, setSorting] = useState<SortingState>([])

  const columns = useMemo<ColumnDef<Contest>[]>(() => [
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
          <Link to="/wallpaper/$id" params={{ id: row.original.id }}>
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
      size: 1,
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
    <StyledTable table={table} />
  )
}

export default ContestPage
