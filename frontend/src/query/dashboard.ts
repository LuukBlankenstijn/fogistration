import { setClientTeam, setTeamClient } from "@/clients/generated-client";
import { listClientsOptions, listClientsQueryKey, listTeamsOptions, listTeamsQueryKey } from "@/clients/generated-client/@tanstack/react-query.gen";
import type { Client, Team } from "@/clients/generated-client/types.gen";
import { useMutation, useQueryClient, useSuspenseQuery } from "@tanstack/react-query"
import { useMemo } from "react"

export type ExtendedClient = Client & { teamId?: string | undefined }

interface TeamNames {
  name: string;
  assigned: boolean;
}

export const useClients = () => {
  const { data: clients } = useSuspenseQuery(listClientsOptions())
  const { data: teams } = useSuspenseQuery(listTeamsOptions())

  const augmentedClients = useMemo((): ExtendedClient[] => {
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
    ...listTeamsOptions(),
    select: (teams): TeamNames[] => teams.map((team) => ({ name: team.name, assigned: !!team.ip }))
  })
}

export const useSetTeamClientMutation = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({ team, clientIp }: { team: Team, clientIp?: string | undefined }) => {
      const client = queryClient.getQueryData<Client[]>(listClientsQueryKey())?.find((c) => c.ip === clientIp)

      const { data: updatedTeam } = await setTeamClient({
        body: {
          clientId: client?.id
        },
        path: {
          id: team.id
        },
        throwOnError: true
      });

      return updatedTeam
    },
    onMutate: ({ team, clientIp }) => {
      const prev = queryClient.getQueryData<Team[]>(listTeamsQueryKey()) ?? []
      const client = queryClient.getQueryData<Client[]>(listClientsQueryKey())?.find((client) => client.ip === clientIp)

      queryClient.setQueryData<Team[]>(
        listTeamsQueryKey(),
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
      queryClient.setQueryData(listTeamsQueryKey(), context?.prev)
    },
    onSettled: async () => {
      await queryClient.invalidateQueries({
        queryKey: listTeamsQueryKey()
      })
    }
  })
}

export const useSetClientTeamMutation = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({ client, teamId }: { client: ExtendedClient, teamId?: string | undefined }) => {
      const { data } = await setClientTeam({
        body: {
          teamId
        },
        path: {
          id: client.id,
        },
        throwOnError: true
      });
      return data;
    },
    onMutate: ({ client, teamId }) => {
      const prev = queryClient.getQueryData<Team[]>(listTeamsQueryKey()) ?? []

      queryClient.setQueryData<Team[]>(
        listTeamsQueryKey(),
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
      queryClient.setQueryData(listTeamsQueryKey(), context?.prev)
    },
    onSettled: async () => {
      await queryClient.invalidateQueries({
        queryKey: listTeamsQueryKey()
      })
    }
  })
}
