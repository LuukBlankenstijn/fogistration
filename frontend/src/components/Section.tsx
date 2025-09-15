import { Suspense, type ReactNode } from "react";
import { Card } from "./Card";

export function Section({ children, fallback }: { children: ReactNode, fallback: string }) {
  return (
    <Card className='m-6 max-w-6xl'>
      <Suspense fallback={<div className="text-sm text-[hsl(var(--muted))]">{fallback}</div>}>
        <section className='space-y-3'>
          {children}
        </section>
      </Suspense>
    </Card>
  )
}
