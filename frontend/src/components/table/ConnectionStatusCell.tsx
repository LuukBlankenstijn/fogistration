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

function formatRelativeDateTime(date: Date): string {
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffHours = diffMs / (1000 * 60 * 60)
  const diffDays = diffHours / 24

  if (diffHours < 24) {
    // Less than 24 hours → show time
    return date.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })
  } else if (diffDays < 7) {
    // Between 1 and 6 days → show number of days
    const days = Math.floor(diffDays)
    return `${days.toString()} day${days !== 1 ? "s" : ""} ago`
  } else {
    // More than a week
    return "more than a week ago"
  }
}
