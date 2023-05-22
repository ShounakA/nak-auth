import { Component, createSignal } from "solid-js";

const Login: Component = () => {
  const urlParams = new URLSearchParams(document.location.search);
  console.log(urlParams)
  console.log(document.URL)
  const initUri = urlParams.get('redirect_uri');
  const [redirectUri, setRedirectUri] = createSignal(initUri);
  return (
   <div class="border-white border-2 w-80 rounded-2xl">
         <div class="py-5 px-5">
           <form class="flex flex-col" action={`/login?redirect_uri=${redirectUri()}`} method="post">
             <div class="flex flex-row justify-between">
             <label for="username">Username: </label>
               <input type="text" id="username" name="username" class="py-1 bg-black text-white" placeholder="something@gmail.com" />
             </div>
             <div class="flex flex-row justify-between">
               <label for="secret">Password: </label>
               <input type="password" id="secret" name="secret" class="py-1 bg-black text-white" />
             </div>
             <div class="flex flex-row justify-end py-2">
               <button class="bg-blue-800 p-1 rounded" value="Login">Authorize</button>
             </div>
           </form>
         </div>
       </div>
  );
};

export default Login;