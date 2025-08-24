import { getWallpaperFile, setWallpaperConfig, setWallpaperFile } from "@/clients/generated-client"
import { getWallpaperConfigOptions, getWallpaperConfigQueryKey, getWallpaperFileQueryKey } from "@/clients/generated-client/@tanstack/react-query.gen"
import type { Layout } from "@/components/WallpaperEditor/types/layout"
import { useMutation, useQueryClient, useSuspenseQuery } from "@tanstack/react-query"

export const useGetWallpaperConfigQuery = (contestId: string) => {
  return useSuspenseQuery(getWallpaperConfigOptions({
    path: {
      contestId
    }
  }))
}

export const useWallpaperConfigMutation = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ layout, contestId }: { layout: Layout, contestId: string }) => {
      return setWallpaperConfig({
        path: {
          contestId
        },
        body: {
          ...layout
        }
      })
    },
    onSuccess: (data, variables) => {
      queryClient.setQueryData(getWallpaperConfigQueryKey({ path: { contestId: variables.contestId } }), data.data)
    },
  })
}

export const useGetWallpaperQuery = (contestId: string) => {
  const options = {
    path: {
      contestId
    }
  }
  return useSuspenseQuery({
    queryFn: async ({ queryKey, signal }) => {
      const { data } = await getWallpaperFile({
        ...options,
        ...queryKey[0],
        signal,
      });
      // TODO: fix this!!!!
      if (!data) {
        return null
      }
      return data as unknown as Blob;
    },
    queryKey: getWallpaperFileQueryKey(options),
  })
}

export const useWallpaperMutation = (contestId: string) => {
  return useMutation({
    mutationFn: async ({ url }: { url: string | undefined }) => {
      const blob = await fetch(url ?? '').then((v) => v.blob())

      const options = {
        path: {
          contestId
        },
        body: url === undefined ? {
          file: blob
        } : undefined
      }
      await setWallpaperFile({
        ...options,
        throwOnError: true
      });
    }
  })
}

