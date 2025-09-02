import { getCurrentUserQueryKey } from "@/clients/generated-client/@tanstack/react-query.gen";
import { MutationCache, QueryCache, QueryClient } from "@tanstack/react-query";

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 30,
    }
  },
  queryCache: new QueryCache({
    onError: (error: unknown) => {
      if (getErrorStatus(error) === BigInt(401)) {
        void queryClient.invalidateQueries({
          queryKey: getCurrentUserQueryKey()
        })
      }
    }
  }),
  mutationCache: new MutationCache({
    onError: (error: unknown) => {
      if (getErrorStatus(error) === BigInt(401)) {
        void queryClient.invalidateQueries({
          queryKey: getCurrentUserQueryKey()
        })
      }
    }
  })
})

function getErrorStatus(error: unknown): bigint | null {
  if (
    typeof error == "object" &&
    error &&
    "status" in error &&
    typeof error.status === "bigint"
  ) return error.status

  return null
}

export default queryClient
