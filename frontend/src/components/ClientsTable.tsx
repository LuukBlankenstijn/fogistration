import { flexRender, getCoreRowModel, getSortedRowModel, useReactTable, type ColumnDef, type SortingState } from "@tanstack/react-table";
import { useMemo, useState } from "react";
import { Dropdown } from "./table/Dropdown";
import { useTeamsQuery } from "@/query/team";
import { useClients, useSetClientTeamMutation, type ExtendedClient } from "@/query/client";
import ConnectionStatusCell from "./ConnectionStatusCell";

export function ClientsTable() {
  const clients = useClients()
  const teams = useTeamsQuery()
  const { mutate } = useSetClientTeamMutation()

  const teamNameById = useMemo(() => {
    const m = new Map<string, string>()
    teams.forEach(t => m.set(t.id, t.name))
    return m
  }, [teams])

  const [sorting, setSorting] = useState<SortingState>([])

  const columns = useMemo<ColumnDef<ExtendedClient>[]>(() => [
    { accessorKey: "id", header: "ID" },
    { accessorKey: "ip", header: "IP" },
    {
      accessorKey: "lastSeen",
      header: "Connected",
      cell: ({ getValue }) => <ConnectionStatusCell lastSeen={getValue<Date>()} />
    },
    {
      id: "team",
      header: "Team",
      accessorFn: (row) => teamNameById.get(row.teamId ?? "") ?? "",
      sortingFn: "alphanumeric",
      cell: ({ row }) => {
        const onSelectTeam = (teamId: string | null) => { mutate({ client: row.original, teamId: teamId ?? undefined }); }

        return (
          <Dropdown
            value={row.original.teamId}
            options={teams}
            valueGenerator={(team) => team.id}
            labelGenerator={(team) => team.name}
            show={(team, value) => !team.ip || value === team.id}
            onChange={onSelectTeam}
            placeholder="No team"
          />
        )
      },
    },
  ], [teams, mutate, teamNameById])

  const table = useReactTable({
    data: clients,
    columns,
    state: { sorting },
    onSortingChange: setSorting,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
  })

  return (
    <div>
      <header className="flex items-baseline justify-between">
        <h2 className="text-lg font-semibold" style={{ color: 'hsl(var(--fg))' }}>Clients</h2>
        <span className="text-xs" style={{ color: 'hsl(var(--muted))' }}>{clients.length} total</span>
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
