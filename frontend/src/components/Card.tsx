import React from "react";

type CardProps = React.ComponentProps<"section">;

export const Card = ({ ref, className = "", ...props }: CardProps & { ref?: React.RefObject<HTMLElement | null> }) => (
  <section
    ref={ref}
    className={`rounded-2xl border border-[hsl(var(--border))] bg-[hsl(var(--panel))] p-6 shadow-sm ${className}`}
    {...props}
  />
);
Card.displayName = "Card";

