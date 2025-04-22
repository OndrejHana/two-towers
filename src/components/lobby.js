import { State, GAME, ERROR, LOADING, MAIN } from "../state";

/**
 * @param {HTMLDivElement} parent
 * @param {State} state
 */
export function renderLobby(parent, clerk, state) {
  const c = new AbortController();
  const html = `<div class="w-full bg-neutral-100 h-full flex md:pt-32 justify-center relative">
    <div class="p-4 bg-neutral-50 shadow-sm rounded-lg h-fit max-w-xl w-full space-y-8">
        <h1 class="text-2xl font-bold text-center">Lobby</h1>
        <div class="bg-white shadow-sm rounded p-2 flex flex-col gap-2">
            <p>sup</p>
            <p>sup</p>
            <p>sup</p>
            <p>sup</p>
        </div>
        <div class="flex justify-between">
            <button type="button" class="p-2 rounded hover:bg-white hover:shadow-sm hover:cursor-pointer text-left" id="leave-game-button">Leave</button>
            <button type="button" class="inline-flex items-center rounded bg-green-600 px-3 py-2 text-sm font-semibold text-white shadow-xs hover:bg-green-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-green-600" id="start-game-button">Ready</button>
        </div>
    </div>
`;

  parent.innerHTML = html;

  const ngb = document.getElementById("start-game-button");
  ngb.addEventListener(
    "mousedown",
    async function (_) {
      let err = state.renderWith(LOADING);
      if (err !== null) {
        state.error = err;
        state.renderWith(ERROR);
      }

      const res = await fetch("game/new", {
        headers: {
          Authorization: `Bearer ${await clerk.session.getToken()}`,
        },
      });
      const payload = await res.json();
      state.context.payload = payload;

      err = state.renderWith(GAME);
      if (err !== null) {
        state.error = err;
        state.renderWith(ERROR);
      }
    },
    {
      signal: c.signal,
    },
  );

  const lgb = document.getElementById("leave-game-button");
  lgb.addEventListener(
    "mousedown",
    function (_) {
      let err = state.renderWith(MAIN);
      if (err !== null) {
        state.error = err;
        state.renderWith(ERROR);
      }
    },
    { signal: c.signal },
  );

  return c;
}
