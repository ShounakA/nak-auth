import { Component } from "solid-js";
import Login from "./Login";

const App: Component = () => {
  return (
   <main class="flex flex-col m-auto bg-black h-screen w-screen text-white items-center justify-center">
    <Login />
   </main>
  );
};

export default App;
