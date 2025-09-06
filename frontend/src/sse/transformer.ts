import type { SseEvent } from "./events"

let modPromise: Promise<typeof import("@/clients/generated-client/transformers.gen")> | null = null
const loadClient = () =>
  (modPromise ??= import("@/clients/generated-client/transformers.gen"))

export async function transformSSE<TData>(
  event: SseEvent,
  data: TData,
): Promise<TData> {
  const mod = await loadClient();
  const exportName = `${event}ResponseTransformer` as keyof typeof mod
  const fn = mod[exportName];

  if (typeof fn !== 'function') {
    return data
  }

  return fn(data) as TData
}
