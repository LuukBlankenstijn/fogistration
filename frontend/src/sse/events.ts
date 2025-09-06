import type { SseResponses } from "@/clients/generated-client";

type EventUnion = SseResponses[200][number];
export type SseEvent = EventUnion["event"];
export type PayloadOf<E extends SseEvent> = Extract<EventUnion, { event: E }>["data"];

export const ALL_EVENTS = [
  "getCurrentUser",
  "getTeam",
  "getClient",
] as const satisfies readonly SseEvent[];
