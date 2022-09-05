import { createRouter } from "./context";
import { z } from "zod";
import { Guid } from 'guid-ts';

export const clientRouter = createRouter()
  .query("getAll", {
    async resolve({ ctx }) {
      return await ctx.prisma.client.findMany();
    },
  })
  .mutation("addOne", {
   input: z.object({id: z.string()}),
    async resolve({ ctx, input }) {
      return await ctx.prisma.client.create({
         data: {
            id: input.id,
            secret: Guid.newGuid().toString()
         }
      })
    }
  })
  .mutation("deleteById", {
    input: z.object( { id: z.string()}),
    async resolve({ ctx, input }) {
      return await ctx.prisma.client.delete({ where: { id: input.id }})
    }
  })
