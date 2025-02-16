import './style.css'
//import { setupCounter } from './counter.ts'

//document.querySelector<HTMLDivElement>('#app')!.innerHTML = `
//  <div>
//    <h1 class="bg-amber-300">Hello world</h1>
//    <a href="https://vite.dev" target="_blank">
//      <img src="${viteLogo}" class="logo" alt="Vite logo" />
//    </a>
//    <a href="https://www.typescriptlang.org/" target="_blank">
//      <img src="${typescriptLogo}" class="logo vanilla" alt="TypeScript logo" />
//    </a>
//    <h1>Vite + TypeScript</h1>
//    <div class="card">
//      <button id="counter" type="button"></button>
//    </div>
//    <p class="read-the-docs">
//      Click on the Vite and TypeScript logos to learn more
//    </p>
//  </div>
//`

//setupCounter(document.querySelector<HTMLButtonElement>('#counter')!)


//async function onStartGame() {
//
//}

const mainMenu = `
    <div class="w-full h-full flex flex-col items-center justify-center gap-2">
        <button class="p-2 rounded hover:bg-neutral-100">start game</button>
        <a class="hover:bg-neutral-100 p-2 rounded" href="http://localhost:8080/auth?provider=google">Sign in with google</a>
    </div>
`;

document.querySelector<HTMLDivElement>("#app")!.innerHTML = mainMenu;
