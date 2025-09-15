import { listUsers, patchUser, putUser, UserRole, type User } from "@/clients/generated-client";
import { listUsersQueryKey, getUserQueryKey } from "@/clients/generated-client/@tanstack/react-query.gen";
import { useMutation, useQueries, useQueryClient, useSuspenseQuery, type QueryKey, type UseQueryOptions } from "@tanstack/react-query";
import { useEffect, useMemo } from "react";

function useUserIdsQuery() {
  const qc = useQueryClient();
  const usersQ = useSuspenseQuery({
    queryKey: listUsersQueryKey(),
    queryFn: async ({ queryKey, signal }) => {
      const { data } = await listUsers({ ...queryKey[0], signal, throwOnError: true });
      return data
    },
    staleTime: 60_000,
    refetchOnWindowFocus: false,
  });

  useEffect(() => {
    usersQ.data.forEach((t) => {
      const key = getUserQueryKey({ path: { id: t.id } } as const);
      qc.setQueryData<User>(key, t);
    });
  }, [qc, usersQ.data]);

  const ids = useMemo(() => usersQ.data.map((u) => u.id), [usersQ.data]);
  // eslint-disable-next-line @tanstack/query/no-rest-destructuring
  return { ...usersQ, data: ids };
}

export function useUsersQuery(): User[] {
  const qc = useQueryClient();
  const { data: ids = [] } = useUserIdsQuery(); // returns bigInt[]

  const queries = useMemo(
    () =>
      ids.map((id) => {
        const key = getUserQueryKey({ path: { id: id } }) as QueryKey;
        return {
          queryKey: key,
          enabled: false,                        // observe-only; SSE updates this key
          staleTime: Infinity,
          notifyOnChangeProps: ['data'] as const,
          placeholderData: () => qc.getQueryData<User>(key), // seed from cache
        } satisfies UseQueryOptions<User, Error, User>;
      }),
    [qc, ids.join(',')]
  );
  return useQueries({
    queries,
    combine: (result) => result.map((r) => r.data as User | undefined).filter((t): t is User => !!t)
  });
}

export function putUserMutation() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async ({ user }: { user: User }) => {
      const { data } = await putUser({
        body: {
          email: user.email,
          username: user.username,
          role: user.role,
        },
        path: {
          id: user.id
        }
      })

      return data
    },
    onMutate: (variables) => {
      const key = getUserQueryKey({ path: { id: variables.user.id } })
      const original = queryClient.getQueryData<User>(key)
      queryClient.setQueryData(key, variables.user)

      return { original, key }
    },
    onError: (_error, _variables, context) => {
      if (context?.original) {
        queryClient.setQueryData(context.key, context.original)
      }
    }

  })
}

export function setUserRoleMutation() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async ({ user, role }: { user: User, role: UserRole }) => {
      const { data } = await patchUser({
        body: {
          role: role

        },
        path: {
          id: user.id
        }
      })

      return data
    },
    onMutate: (variables) => {
      const key = getUserQueryKey({ path: { id: variables.user.id } })
      const old = queryClient.getQueryData<User>(key)
      if (!old) return

      queryClient.setQueryData<User>(key, (old) => {
        if (!old) {
          return {
            ...variables.user,
            role: variables.role
          }
        }
        return {
          ...old,
          role: variables.role,
        }
      })

      return { old }
    },
    onError: (_error, variables, context) => {
      if (!context) return
      const key = getUserQueryKey({ path: { id: variables.user.id } })
      queryClient.setQueryData<User>(key, context.old)
    }
  })
}
