import { type User } from "@/clients/generated-client";
import { devLoginMutation, getCurrentUserOptions, getCurrentUserQueryKey, loginMutation, logoutMutation, oidcEnabledOptions } from "@/clients/generated-client/@tanstack/react-query.gen";
import { useMutation, useQuery, useQueryClient, useSuspenseQuery } from "@tanstack/react-query";

export const useLogin = () => {
  const queryClient = useQueryClient()
  return useMutation({
    ...loginMutation(),
    onSuccess: (data) => {
      queryClient.setQueryData<User>(getCurrentUserQueryKey(), data)
    },
    onError: (error) => {
      if (error.status === BigInt(403)) {
        queryClient.removeQueries({
          queryKey: getCurrentUserQueryKey()
        })
      }
    }
  })
}

export const useGetCurrentUser = () => {
  return useSuspenseQuery({
    ...getCurrentUserOptions()
  })
}

export const useLogout = () => {
  const queryClient = useQueryClient()
  return useMutation({
    ...logoutMutation(),
    onSuccess: () => {
      void queryClient.invalidateQueries()
    }
  })
}

export const useDevLogin = () => {
  const queryClient = useQueryClient();
  return useMutation({
    ...devLoginMutation(),
    onSuccess: (data) => {
      queryClient.setQueryData<User>(getCurrentUserQueryKey(), data)
    },
    onError: (error) => {
      if (error.status === BigInt(403)) {
        queryClient.removeQueries({
          queryKey: getCurrentUserQueryKey()
        })
      }
    }
  });
};

export const useOIDCEnabledCheck = () => {
  return useQuery({
    ...oidcEnabledOptions(),
    retry: false
  })
}
