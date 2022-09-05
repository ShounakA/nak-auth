// src/server/router/index.ts
import { createRouter } from "./context";
import superjson from "superjson";

import { clientRouter } from "./client";

export const appRouter = createRouter()
  .transformer(superjson)
  .merge("client.", clientRouter);

// export type definition of API
export type AppRouter = typeof appRouter;
