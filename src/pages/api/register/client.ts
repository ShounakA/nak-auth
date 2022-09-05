// src/pages/api/register/client
import { Guid } from "guid-ts";
import type { NextApiRequest, NextApiResponse } from "next";
import { prisma } from "../../../server/db/client";

//TODO remove this should only be able to create clients from the client
const clients = async (req: NextApiRequest, res: NextApiResponse) => {
   switch (req.method) {
      // case 'POST':
      //    const name: string = req.body.name;
      //    const item = await prisma.client.create({
      //       data: {
      //          name: name,
      //          secret: Guid.newGuid().toString()
      //       }
      //    })
      //    res.status(201).json(item);
      // case 'GET':
      //    if (req.query['name']){
      //       const name: string = req.query['name'].toString();
      //       const client = await prisma.client.findFirst({ where : { name : name }});
      //       res.status(200).json(client)
      //    } else{
      //       const clients = await prisma.client.findMany()
      //       res.status(200).json(clients);
      //    }
      // case 'DELETE':{
      //    if (req.query['id']){
      //       const name: string = req.query['id'].toString();
      //       const client = await prisma.client.deleteMany({ where : { name : name }});
      //       res.status(200).json(client)
      //    }
      //    const client = await prisma.client.deleteMany({ where : { name : "" }});
      //    res.status(200).json(client)
      // }
      default:
         res.status(400).json({});
   }

};

export default clients;