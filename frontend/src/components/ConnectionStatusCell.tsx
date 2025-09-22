import { formatRelativeDateTime } from "@/utils/date";
import { useEffect, useState } from "react";

function ConnectionStatusCell({ lastSeen }: { lastSeen: Date }) {
  const [now, setNow] = useState(() => Date.now());

  useEffect(() => {
    const id = setInterval(() => { setNow(Date.now()); }, 1000);
    return () => { clearInterval(id); };
  }, []);

  const lastSeenMs = new Date(lastSeen).getTime();
  const diff = now - lastSeenMs;

  if (diff <= 5000) {
    return <span className="text-green-500">Connected</span>;
  }

  return (
    <span className="text-red-500">
      Disconnected (last seen {formatRelativeDateTime(lastSeen)})
    </span>
  );
}

export default ConnectionStatusCell
