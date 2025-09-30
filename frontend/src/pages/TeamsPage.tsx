import { getCoreRowModel, getSortedRowModel, useReactTable, type ColumnDef, type SortingState } from "@tanstack/react-table";
import { useMemo, useState } from "react";
import type { Team } from "@/clients/generated-client";
import { useSetTeamClientMutation, useTeamsQuery } from "@/query/team";
import { useClients } from "@/query/client";
import { Dropdown } from "@/components/table/Dropdown";
import StyledTable from "@/components/table/Table";

const TeamsPage = () => {
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
    <StyledTable table={table} />
  )
}

export default TeamsPage
