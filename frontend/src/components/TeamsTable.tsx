import { flexRender, getCoreRowModel, getSortedRowModel, useReactTable, type ColumnDef, type SortingState } from "@tanstack/react-table";
import { useMemo, useState } from "react";
import { Dropdown } from "./table/Dropdown";
import type { Team } from "@/clients/generated-client";
import { useSetTeamClientMutation, useTeamsQuery } from "@/query/team";
import { useClients } from "@/query/client";

export function TeamsTable() {
  const clients = useClients()
  const teams = useTeamsQuery()
  const { mutate } = useSetTeamClientMutation()

  const availableIps = useMemo(() => {
    return clients.filter((client) => !!client.teamId).map((client) => client.ip)
  }, [clients])

  const [sorting, setSorting] = useState<SortingState>([])

  const columns = useMemo<ColumnDef<Team>[]>(() => [
    { accessorKey: "id", header: "ID" },
    { accessorKey: "name", header: "Name" },
    {
      id: "client",
      header: "Client",
      accessorKey: "ip",
      cell: ({ row }) => {
        const onSelectClient = (clientIp: string | null) => {
          mutate({ team: row.original, clientIp: clientIp ?? undefined })
        }

        return (
          <Dropdown
            value={row.original.ip}
            options={clients}
            valueGenerator={(client) => client.ip}
            labelGenerator={(client) => client.ip}
            show={(client, value) => !availableIps.includes(client.ip) || value === client.ip}
            onChange={onSelectClient}
            placeholder="No client"
          />
        )
      },
    },
  ], [clients, mutate, availableIps])

  const table = useReactTable({
    data: teams,
    columns,
    state: { sorting },
    onSortingChange: setSorting,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
  })

  return (
    <div>
      <header className="flex items-baseline justify-between">
        <h2 className="text-lg font-semibold" style={{ color: 'hsl(var(--fg))' }}>Teams</h2>
        <span className="text-xs" style={{ color: 'hsl(var(--muted))' }}>{teams.length} total</span>
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
