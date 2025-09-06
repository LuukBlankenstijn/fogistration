import type { SseEvent } from "./events";

export type SseHandler = (e: MessageEvent<string>) => void;

class SSEBus {
  private es: EventSource | null = null;
  private handlers = new Map<SseEvent | "message", Set<SseHandler>>();
  private native = new Map<SseEvent, EventListener>();
  constructor(private url = "/api/sse", private creds = false) { this.connect(); }

  private connect() {
    this.es?.close();
    this.es = new EventSource(this.url, { withCredentials: this.creds });
    this.es.onerror = () => { /* optional: backoff & reconnect */ };
  }

  add(event: SseEvent | "message", fn: SseHandler): () => void {
    let set = this.handlers.get(event);
    if (!set) {
      set = new Set(); this.handlers.set(event, set);
      if (event !== "message" && this.es && !this.native.has(event)) {
        const cb: EventListener = (e) => { this.emit(event, e as MessageEvent<string>); };
        this.es.addEventListener(event, cb);
        this.native.set(event, cb);
      }
    }
    set.add(fn);
    return () => { this.remove(event, fn); };
  }

  private remove(event: SseEvent | "message", fn: SseHandler) {
    const set = this.handlers.get(event);
    if (!set) return;
    set.delete(fn);
    if (set.size === 0) {
      this.handlers.delete(event);
      if (event !== "message") {
        const cb = this.native.get(event);
        if (cb && this.es) this.es.removeEventListener(event, cb);
        this.native.delete(event);
      }
    }
  }

  private emit(event: SseEvent | "message", e: MessageEvent<string>) {
    const set = this.handlers.get(event);
    if (!set) return;
    for (const fn of set) fn(e);
  }
}

let singleton: SSEBus | null = null;
export function getSSE(url?: string, withCredentials?: boolean) {
  singleton ??= new SSEBus(url, !!withCredentials);
  return singleton;
}
