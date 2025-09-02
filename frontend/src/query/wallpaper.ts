import { deleteWallpaperFile, deleteWallpaperLayout, getWallpaperFile, getWallpaperLayout, putWallpaperFile, putWallpaperLayout, type WallpaperLayout } from "@/clients/generated-client"
import { getWallpaperFileQueryKey, getWallpaperLayoutQueryKey } from "@/clients/generated-client/@tanstack/react-query.gen"
import { DEFAULT_WALLPAPER_LAYOUT } from "@/components/WallpaperEditor/types"
import { useMutation, useQueryClient, useSuspenseQuery } from "@tanstack/react-query"
import { useNavigate } from "@tanstack/react-router"

export const useGetWallpaperConfigQuery = (id: number) => {
  const queryClient = useQueryClient()
  return useSuspenseQuery({
    queryFn: async ({ queryKey, signal }) => {
      const { data, error } = await getWallpaperLayout({
        ...queryKey[0],
        signal,
      });
      if (error !== undefined) {
        queryClient.setQueryData(queryKey, DEFAULT_WALLPAPER_LAYOUT)
        throw error
      }
      return data;
    },
    queryKey: getWallpaperLayoutQueryKey({ path: { id } }),
    retry: false,

  })
}

export const useGetWallpaperQuery = (id: number) => {
  const queryClient = useQueryClient()
  return useSuspenseQuery({
    queryFn: async ({ queryKey, signal }) => {
      const { data, error } = await getWallpaperFile({
        ...queryKey[0],
        signal,
      });
      if (error !== undefined) {
        queryClient.setQueryData(queryKey, null)
        throw error
      }
      return data;
    },
    queryKey: getWallpaperFileQueryKey({ path: { id } }),
    retry: false
  })
}

export const useWallpaperMutation = (id: number) => {
  const queryClient = useQueryClient()
  const options = {
    path: {
      id
    }
  }
  const navigate = useNavigate()
  return useMutation({
    mutationFn: async ({
      layout,
      file,
    }: { layout: WallpaperLayout | undefined; file: File | Blob | null }) => {
      const layoutReq = (layout
        ? putWallpaperLayout({
          path: { id },
          body: layout,
          throwOnError: true,
        })
        : deleteWallpaperLayout({
          path: { id },
          throwOnError: true,
        })
      );

      const fileReq = (file
        ? putWallpaperFile({
          ...options,
          body: file,
          throwOnError: true,
          bodySerializer: (b: Blob) => b,
        })
        : deleteWallpaperFile({
          ...options,
          throwOnError: true,
        })
      );

      const [{ data: newLayout }, { data: newFile }] = await Promise.all([layoutReq, fileReq]);

      // need to change void to null, otherwise in onSuccess data will be an empty object
      // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
      return { newLayout, newFile: newFile ?? null };
    },
    onSuccess: (data) => {
      if (data.newLayout) {
        queryClient.setQueryData(getWallpaperLayoutQueryKey(options), data.newLayout)
      }
      if (data.newFile instanceof Blob) {
        queryClient.setQueryData(getWallpaperFileQueryKey(options), data.newFile)
      } else {
        queryClient.setQueryData(getWallpaperFileQueryKey(options), null)
      }
      void navigate({ to: "/dashboard" })
    }
  })
}

