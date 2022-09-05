import type { NextApiRequest, NextApiResponse } from "next";
import { prisma } from "../../server/db/client";
import { sign, SignOptions, Secret } from "jsonwebtoken";

const auth = async (req: NextApiRequest, res: NextApiResponse) => {
   switch (req.method) {
      case 'POST':
         if (!req.body.client_id) {
            res.status(400).json("Must specificy 'client_id'");
         }
         if (!req.body.client_secret) {
            res.status(400).json("Must specificy 'client_secret'");
         }
         if (!req.body.grant_type) {
            res.status(400).json("Must specificy 'grant_type'");
         }
         const id: string = req.body.client_id;
         const secret: string= req.body.client_secret;
         const grant_type: string = req.body.grant_type;
         switch (grant_type) {
            case 'client_credentials':
               const client = await prisma.client.findUnique({
                  where: {
                     id
                  }
               });
               if (!client || client.secret !== secret) {
                  res.status(401).json("Unknown client credentials");
                  break;
               }
               const scope = ["all"];
               const token = generate_token(scope, id);
               res.status(200).json({ "scope": scope, "token": token});
            default:
               res.status(400).json("Invalid grant type.")
         }
      default:
         res.status(400).json("Unknown request");
   }

};


const generate_token = (scope: string[], client_id: string) => {
   const payload = {
      id: client_id,
      access_types : scope,
   }

   const privateKey = {
      key: process.env['PRIVATE_KEY'],
      passphrase: process.env['PASSPHRASE']
   } as Secret;

   if (!privateKey) {
      throw Error('No private key defined');
   }
   const signInOptions: SignOptions = {
      // RS256 uses a public/private key pair. The API provides the private key
      // to generate the JWT. The client gets a public key to validate the
      // signature
      algorithm: 'RS256',
      expiresIn: '1h',
    };
  
    return sign(payload, privateKey, signInOptions);
}

export default auth;