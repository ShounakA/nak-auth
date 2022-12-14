import type { NextPage } from "next";
import Head from "next/head";
import { useState } from "react";
import { trpc } from "../utils/trpc";

const getClients = () => {
   const clients = trpc.useQuery(['client.getAll']).data;
   const list = []
   if (clients){
      for (const client of clients) {
         list.push(
          <li key={client.id}>
           { client.id } - { client.secret }
          </li>)
      }
   } else {
      list.push(<li> No Clients Registered </li>)
   }
   return list;
}

const Clients: NextPage = () => {
  // eslint-disable-next-line prefer-const
  const [client, setClient] = useState('');
  const [del_client, setDelClient] = useState('');
  const add = trpc.useMutation('client.addOne')
  const del = trpc.useMutation('client.deleteById');
  const handleSubmit = () => {
    add.mutate({ id: client });
  }
  const handleDelSubmit = () => {
    del.mutate({ id: del_client });
  }
  return (
    <>
      <Head>
        <title>nak Auth</title>
        <meta name="description" content="Generated by create-t3-app" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <main className="container mx-auto flex flex-col items-center justify-center min-h-screen p-4">
        <h1 className="text-5xl md:text-[5rem] leading-normal font-extrabold text-gray-700">
          nak <span className="text-purple-300">Clients</span> 
        </h1>
        <div className="flex flex-row">
          <div className="flex flex-col">
            <form onSubmit={handleSubmit}>
              <label>
                Add New: 
                <input 
                  type="text" 
                  onChange={e => setClient(e.target.value)} />
              </label>
              <input type="submit" value="Add" />
            </form>
            <form onSubmit={handleDelSubmit}>
              <label>
                Delete: 
                <input 
                  type="text" 
                  onChange={e => setDelClient(e.target.value)} />
              </label>
              <input type="submit" value="Delete" />
            </form>
          </div>
          <div className="w-10 mx-auto justify-center"></div>
          <div>
            <h3 className="text-2xl leading-normal font-extrabold text-gray-700"> Registered Clients </h3>
            <ul className="text-gray-700">
            { getClients() }
            </ul>
          </div>
        </div>
      </main>
    </>
  );
};

export default Clients;

