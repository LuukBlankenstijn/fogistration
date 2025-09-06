import { useEffect, useRef } from "react";
import { getSSE } from "./sse";
import { ALL_EVENTS, type SseEvent, type PayloadOf } from "./events";

type MaybePromise<T> = T | Promise<T>;
type AllHandler = <E extends SseEvent>(
  event: E,
  data: PayloadOf<E>,
  raw: MessageEvent<string>,
) => MaybePromise<void>;

export function useSSEAll(onEvent: AllHandler) {
  const ref = useRef(onEvent);
  useEffect(() => { ref.current = onEvent; }, [onEvent]);

  useEffect(() => {
    const bus = getSSE();

    const unsubs = ALL_EVENTS.map((evt) =>
      bus.add(evt, (e) => {
        // parse once per event, swallow promise so listener is effectively void
        try {
          const data = JSON.parse(e.data) as PayloadOf<typeof evt>;
          void Promise.resolve(ref.current(evt, data, e));
        } catch {
          // optional: console.warn("[SSE] bad JSON for", evt, e.data);
        }
      })
    );

    return () => { unsubs.forEach((u) => { u(); }); };
  }, []);
}
