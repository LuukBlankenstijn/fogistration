import type { User } from "@/clients/generated-client"
import { putUserMutation, useUsersQuery } from "@/query/users"
import { flexRender, getCoreRowModel, getSortedRowModel, useReactTable, type ColumnDef, type SortingState } from "@tanstack/react-table"
import { useCallback, useMemo, useState } from "react"
import { EditUserModal } from "./EditUserModal"

export const UsersTable = () => {
  const users = useUsersQuery()
  const { mutateAsync } = putUserMutation()

  const [sorting, setSorting] = useState<SortingState>([])
  const [editing, setEditing] = useState<User | null>(null)
  const [open, setOpen] = useState(false)

  const columns = useMemo<ColumnDef<User>[]>(() => [
    { accessorKey: "id", header: "ID" },
    { accessorKey: "username", header: "Name" },
    { accessorKey: "email", header: "Email" },
    { accessorKey: "role", header: "Role" },
    {
      id: "actions",
      header: "",
      cell: ({ row }) => (
        <button
          type="button"
          className="rounded-md border px-2 py-1 text-xs"
          style={{ borderColor: "hsl(var(--border))", color: "hsl(var(--fg))" }}
          onClick={() => { setEditing(row.original); setOpen(true) }}>
          Edit
        </button>
      ),
      enableSorting: false,
    },
  ], [])


  const table = useReactTable({
    data: users,
    columns,
    state: { sorting },
    onSortingChange: setSorting,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
  })

  const onSave = useCallback(async (u: User) => {
    await mutateAsync({ user: u })
    setOpen(false)
    setEditing(null)
  }, [mutateAsync])

  return (
    <div>
      <header className="flex items-baseline justify-between">
        <h2 className="text-lg font-semibold" style={{ color: 'hsl(var(--fg))' }}>Users</h2>
        <span className="text-xs" style={{ color: 'hsl(var(--muted))' }}>{users.length} total</span>
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
                  No users found.
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>


      <EditUserModal
        open={open}
        user={editing}
        onClose={() => { setOpen(false); setEditing(null) }}
        onSave={onSave}
      />
    </div>
  )
}
