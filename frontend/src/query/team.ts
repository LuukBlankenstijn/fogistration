import { listTeams, setTeamClient, type Client, type Team } from "@/clients/generated-client";
import { getTeamQueryKey, listClientsQueryKey, listTeamsQueryKey } from "@/clients/generated-client/@tanstack/react-query.gen"
import { useMutation, useQueries, useQueryClient, useSuspenseQuery, type QueryKey, type UseQueryOptions } from "@tanstack/react-query"
import { useEffect, useMemo } from "react";


interface TeamNames {
  name: string;
  assigned: boolean;
}


function useTeamIdsQuery() {
  const qc = useQueryClient();
  const teamsQ = useSuspenseQuery({
    queryKey: listTeamsQueryKey(),
    queryFn: async ({ queryKey, signal }) => {
      const { data } = await listTeams({ ...queryKey[0], signal, throwOnError: true });
      return data
    },
    staleTime: 60_000,
    refetchOnWindowFocus: false,
  });

  useEffect(() => {
    teamsQ.data.forEach((t) => {
      const key = getTeamQueryKey({ path: { id: t.id } } as const);
      qc.setQueryData<Team>(key, t);
    });
  }, [qc, teamsQ.data]);

  const ids = useMemo(() => teamsQ.data.map((t) => t.id), [teamsQ.data]);
  // eslint-disable-next-line @tanstack/query/no-rest-destructuring
  return { ...teamsQ, data: ids };
}

export function useTeamsQuery(): Team[] {
  const qc = useQueryClient();
  const { data: ids = [] } = useTeamIdsQuery(); // returns number[]

  const queries = useMemo(
    () =>
      ids.map((id) => {
        const key = getTeamQueryKey({ path: { id } }) as QueryKey;
        return {
          queryKey: key,
          enabled: false,                        // observe-only; SSE updates this key
          staleTime: Infinity,
          notifyOnChangeProps: ['data'] as const,
          placeholderData: () => qc.getQueryData<Team>(key), // seed from cache
        } satisfies UseQueryOptions<Team, Error, Team>;;
      }),
    [qc, ids.join(',')]
  );
  return useQueries({
    queries,
    combine: (result) => result.map((r) => r.data as Team | undefined).filter((t): t is Team => !!t)
  });
}


export const useTeamNames = (): TeamNames[] => {
  const teams = useTeamsQuery()
  return teams.map((team): TeamNames => ({ name: team.name, assigned: !!team.ip }))
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
      const queryKey = getTeamQueryKey({ path: { id: team.id } })
      const prev = queryClient.getQueryData<Team>(queryKey)
      queryClient.setQueryData<Team>(
        queryKey,
        (old) => {
          if (!old) return old
          return {
            ...old,
            ip: clientIp,
          }
        }
      )

      return { prev, queryKey }
    },
    onError: (_error, _variable, context) => {
      console.error(_error)
      if (!context) return
      queryClient.setQueryData(context.queryKey, context.prev)
    },
  })
}
