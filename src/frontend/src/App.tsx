import { Component } from "solid-js";

const App: Component = () => {
  return (
   <main class="flex flex-col m-auto bg-black h-screen w-screen text-white items-center justify-center">
    <div class="border-white border-2 w-80 rounded-2xl">
      <div class="py-5 px-5">
        <form class="flex flex-col" onSubmit={() => console.log('test')}>
          <div class="flex flex-row justify-between">
          <label for="username">Username:</label>
            <input type="text" id="username" class="py-1 bg-black text-white" placeholder="something@gmail.com" />
          </div>
          <div class="flex flex-row justify-between">
            <label for="secret">Secret:</label>
            <input type="password" id="secret" class="py-1 bg-black text-white" />
          </div>
          <div class="flex flex-row justify-end py-2">
            <input class="bg-blue-800 p-1 rounded" value="Login" type="submit" />
          </div>
        </form>
      </div>
    </div>
   </main>
  );
};

export default App;
