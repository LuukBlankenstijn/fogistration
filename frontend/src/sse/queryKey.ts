import type { SseEvent } from "./events";

let modPromise: Promise<typeof import("@/clients/generated-client/@tanstack/react-query.gen")> | null = null;
const loadClient = () =>
  (modPromise ??= import("@/clients/generated-client/@tanstack/react-query.gen")); // cached

export async function queryKeyFor(
  event: SseEvent,
  opts?: unknown
): Promise<unknown[] | null> {
  return queryKey(event, opts);
}

async function queryKey(
  event: string,
  opts?: unknown
): Promise<unknown[] | null> {
  const mod = await loadClient();
  const exportName = `${event}QueryKey` as keyof typeof mod;
  const fn = mod[exportName];

  if (typeof fn !== "function") {
    return null;
  }
  return (fn as (o?: unknown) => unknown[])(opts);
}

export async function queryKeyForList(
  event: SseEvent,
): Promise<unknown[] | null> {
  const listOperationId = event.replace(/^get(\w+)/, (_m: string, word: string) => {
    return "list" + word + "s";
  });

  return queryKey(listOperationId)
}
