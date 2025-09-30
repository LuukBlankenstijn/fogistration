import StyledTable from "@/components/table/Table";
import ConnectionStatusCell from "@/components/table/ConnectionStatusCell"
import { Dropdown } from "@/components/table/Dropdown"
import { useClients, useSetClientTeamMutation, type ExtendedClient } from "@/query/client"
import { useTeamsQuery } from "@/query/team"
import { type SortingState, type ColumnDef, useReactTable, getCoreRowModel, getSortedRowModel } from "@tanstack/react-table"
import { useMemo, useState } from "react"

const ClientsPage = () => {
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

  return <StyledTable tableData={table} />
}

export default ClientsPage
