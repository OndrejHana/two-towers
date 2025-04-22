import { Clerk } from "@clerk/clerk-js";
import { State, LOBBY, ERROR } from "../state";

/**
 * @param {HTMLDivElement} parent
 * @param {Clerk} clerk
 * @param {State} state
 */
export function renderMain(parent, clerk, state) {
  const c = new AbortController();

  const html = `<div class="w-full bg-neutral-100 h-full flex md:pt-32 justify-center relative">
    <div class="p-4 bg-neutral-50 shadow-sm rounded-lg h-fit max-w-64 w-full space-y-8">
        <h1 class="text-2xl font-bold text-center">Two towers</h1>
        <div class="space-y-4 text-sm font-semibold">
            <button type="button" class="inline-flex items-center rounded bg-green-600 px-3 py-2 text-sm font-semibold text-white shadow-xs hover:bg-green-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-green-600 w-full" id="new-game-button">New Game</button>
            <div class="flex flex-col hadow-sm ">
                <input type="text" class="p-1 border border-neutral-200 outline-none block grow rounded-t font-mono font-thin shadow-sm"/>
                <button type="button" class="inline-flex items-center rounded-b bg-green-600 px-3 py-2 text-sm font-semibold text-white shadow-xs hover:bg-green-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-green-600" id="start-game-button">Join lobby</button>
            </div>
        </div>
    </div>
    <div class="absolute top-0 right-0 p-4" ><div id="user-button-container"></div></div>
</div>`;

  parent.innerHTML = html;

  const ubc = document.getElementById("user-button-container");
  clerk.mountUserButton(ubc);

  const ngb = document.getElementById("new-game-button");
  ngb.addEventListener(
    "mousedown",
    async function (e) {
      const err = state.renderWith(LOBBY);
      if (err !== null) {
        state.error = err;
        state.renderWith(ERROR);
      }
      //const res = await fetch("game/new", {
      //  headers: {
      //    Authorization: `Bearer ${await clerk.session.getToken()}`,
      //  },
      //});
      //const payload = await res.json();
      //state;
    },
    {
      signal: c.signal,
    },
  );

  return c;
}
