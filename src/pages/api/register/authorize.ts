import type { NextApiRequest, NextApiResponse } from "next";

const auth = async (req: NextApiRequest, res: NextApiResponse) => {
   switch (req.method) {
      case 'POST':

      case 'GET':

      case 'DELETE':

      default:
         res.status(400).json({});
   }

};

export default auth;