import { setTeam, setTeamClient, type ModelsClient, type ModelsTeam } from "@/clients/generated-client"
import { getAllClientsOptions, getAllClientsQueryKey, getAllTeamsOptions, getAllTeamsQueryKey } from "@/clients/generated-client/@tanstack/react-query.gen"
import { useMutation, useQueryClient, useSuspenseQuery } from "@tanstack/react-query"
import { useMemo } from "react"

export type Client = ModelsClient & { teamId?: string | undefined }

interface TeamNames {
  name: string;
  assigned: boolean;
}

export const useClients = () => {
  const { data: clients } = useSuspenseQuery(getAllClientsOptions())
  const { data: teams } = useSuspenseQuery(getAllTeamsOptions())

  const augmentedClients = useMemo((): Client[] => {
    const teamIdByIp = new Map<string, string>()
    teams.forEach(team => {
      if (team.ip) {
        teamIdByIp.set(team.ip, team.id)
      }
    });

    return clients.map(client => ({
      ...client,
      teamId: teamIdByIp.get(client.ip)
    }))
  }, [clients, teams])


  return augmentedClients
}

export const useTeamNames = () => {
  return useSuspenseQuery({
    ...getAllTeamsOptions(),
    select: (teams): TeamNames[] => teams.map((team) => ({ name: team.name, assigned: !!team.ip }))
  })
}

export const useSetTeamClientMutation = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({ team, clientIp }: { team: ModelsTeam, clientIp?: string | undefined }) => {
      const client = queryClient.getQueryData<ModelsClient[]>(getAllClientsQueryKey())?.find((c) => c.ip === clientIp)

      const { data: updatedTeam } = await setTeamClient({
        body: {
          clientId: client?.id
        },
        path: {
          teamId: team.id
        },
        throwOnError: true
      });

      return updatedTeam
    },
    onMutate: ({ team, clientIp }) => {
      const prev = queryClient.getQueryData<ModelsTeam[]>(getAllTeamsQueryKey()) ?? []
      const client = queryClient.getQueryData<ModelsClient[]>(getAllClientsQueryKey())?.find((client) => client.ip === clientIp)

      queryClient.setQueryData<ModelsTeam[]>(
        getAllTeamsQueryKey(),
        (old = []) =>
          old.map(t => {
            if (t.id === team.id) {
              return { ...t, ip: client?.ip }
            }
            return t
          })
      )

      return { prev }
    },
    onError: (_error, _variable, context) => {
      console.error(_error)
      queryClient.setQueryData(getAllTeamsQueryKey(), context?.prev)
    },
    onSettled: async () => {
      await queryClient.invalidateQueries({
        queryKey: getAllTeamsQueryKey()
      })
    }
  })
}

export const useSetClientTeamMutation = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({ client, teamId }: { client: Client, teamId?: string | undefined }) => {
      const { data } = await setTeam({
        path: {
          clientId: client.id,
        },
        body: {
          teamId
        },
        throwOnError: true
      });
      return data;
    },
    onMutate: ({ client, teamId }) => {
      const prev = queryClient.getQueryData<ModelsTeam[]>(getAllTeamsQueryKey()) ?? []

      queryClient.setQueryData<ModelsTeam[]>(
        getAllTeamsQueryKey(),
        (old = []) =>
          old.map(t => {
            if (t.id === client.teamId) {
              return { ...t, ip: undefined }
            }
            if (t.id === teamId) {
              return { ...t, ip: client.ip }
            }
            return t
          })
      )

      return { prev }
    },
    onError: (_error, _variable, context) => {
      queryClient.setQueryData(getAllTeamsQueryKey(), context?.prev)
    },
    onSettled: async () => {
      await queryClient.invalidateQueries({
        queryKey: getAllTeamsQueryKey()
      })
    }
  })
}
