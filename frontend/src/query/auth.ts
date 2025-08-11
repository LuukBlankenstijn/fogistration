import type { GetCurrentUserResponse } from "@/clients/generated-client";
import { getCurrentUserQueryKey, loginDevMutation, loginMutation, logoutMutation } from "@/clients/generated-client/@tanstack/react-query.gen";
import { useMutation, useQueryClient } from "@tanstack/react-query";

export const useLogin = () => {
  const queryClient = useQueryClient()
  return useMutation({
    ...loginMutation(),
    onSuccess: (data) => {
      queryClient.setQueryData(getCurrentUserQueryKey(), (): GetCurrentUserResponse => ({
        user: data,
        authenticated: true
      }))
    }
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
    ...loginDevMutation(),
    onSuccess: (data) => {
      queryClient.setQueryData(getCurrentUserQueryKey(), (): GetCurrentUserResponse => ({
        user: data,
        authenticated: true
      }))
    },
  });
};

