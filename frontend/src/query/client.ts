import { listClients, setClientTeam, type Client, type Team } from "@/clients/generated-client";
import { getClientQueryKey, getTeamQueryKey, listClientsQueryKey } from "@/clients/generated-client/@tanstack/react-query.gen";
import { useMutation, useQueries, useQueryClient, useSuspenseQuery, type UseQueryOptions } from "@tanstack/react-query";
import { useEffect, useMemo } from "react";
import { useTeamsQuery } from "./team";

export type ExtendedClient = Client & { teamId?: string | undefined }

function useClientIdsQuery() {
  const qc = useQueryClient();
  const clientsQ = useSuspenseQuery({
    queryKey: listClientsQueryKey(),
    queryFn: async ({ queryKey, signal }) => {
      const { data } = await listClients({ ...queryKey[0], signal, throwOnError: true });
      return data
    },
    staleTime: 60_000,
    refetchOnWindowFocus: false,
  });

  useEffect(() => {
    clientsQ.data.forEach((t) => {
      const key = getClientQueryKey({ path: { id: t.id } } as const);
      qc.setQueryData<Client>(key, t);
    });
  }, [qc, clientsQ.data]);

  const ids = useMemo(() => clientsQ.data.map((c) => c.id), [clientsQ.data]);
  // eslint-disable-next-line @tanstack/query/no-rest-destructuring
  return { ...clientsQ, data: ids };
}

export function useClientsQuery(): Client[] {
  const qc = useQueryClient()
  const { data: ids } = useClientIdsQuery()

  const queries = useMemo(
    () =>
      ids.map((id) => {
        const key = getClientQueryKey({ path: { id } })
        return {
          queryKey: key,
          enabled: false,
          staleTime: Infinity,
          notifyOnChangeProps: ["data"] as const,
          placeholderData: () => qc.getQueryData<Client>(key)
        } satisfies UseQueryOptions<Client, Error, Client>
      }),
    [qc, ids.join(",")]
  )

  return useQueries({
    queries,
    combine: (result) => result.map((r) => r.data as Client | undefined).filter((c): c is Client => !!c)
  })
}

export const useClients = () => {
  const clients = useClientsQuery()
  const teams: Team[] = useTeamsQuery()

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
      const prev = queryClient.getQueryData<Team>(getTeamQueryKey({ path: { id: (client.teamId ?? "") } })) ?? undefined

      if (teamId) {
        queryClient.setQueryData<Team>(
          getTeamQueryKey({ path: { id: teamId } }),
          (old) => {
            if (!old) return old
            return {
              ...old,
              ip: client.ip,
            }
          }
        )
      } else if (client.teamId) {
        queryClient.setQueryData<Team>(
          getTeamQueryKey({ path: { id: client.teamId } }),
          (old) => {
            if (!old) return old
            return {
              ...old,
              ip: undefined,
            }
          }
        )
      }

      return { prev }
    },
    onError: (_error, variable, context) => {
      if (variable.teamId) {
        queryClient.setQueryData(getTeamQueryKey({ path: { id: variable.teamId } }), context?.prev)
      } else if (variable.client.teamId) {
        queryClient.setQueryData(getTeamQueryKey({ path: { id: variable.client.teamId } }), context?.prev)
      }
    },
  })
}
