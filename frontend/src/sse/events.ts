import type { SseResponses } from "@/clients/generated-client";

type EventUnion = SseResponses[200][number];
export type SseEvent = EventUnion["event"];
export type PayloadOf<E extends SseEvent> = Extract<EventUnion, { event: E }>["data"];

export const ALL_EVENTS = [
  "getUser",
  "getTeam",
  "getClient",
] as const satisfies readonly SseEvent[];
