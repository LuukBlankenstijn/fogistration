import type { LinkProps, RegisteredRouter } from "@tanstack/react-router";

export type Route = NonNullable<LinkProps<RegisteredRouter>['to']>
