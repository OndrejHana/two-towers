import { Clerk } from "@clerk/clerk-js";
import "./style.css";

import { html, render } from "lit";

const clerkPubKey = import.meta.env.VITE_CLERK_PUBLISHABLE_KEY;

const APP_STATE_LOADING = 1;
const APP_STATE_LOGIN = 2;
const APP_STATE_MAIN = 3;
const APP_STATE_LOBBY = 4;
const APP_STATE_GAME = 5;
const APP_STATE_ERR = 6;

const $ = {
  app: () => document.getElementById("app"),
  userBtn: () => document.getElementById("user-button"),
  newGameBtn: () => document.getElementById("new-game-button"),
  lobbyPlayerList: () => document.getElementById("lobby-player-list"),
  connectionString: () => document.getElementById("connection-string"),
  clerk: null,
  cancelEffects: null,
};

function renderLogin(clerk) {
  const template = html`
    <div class="w-full h-full flex items-center justify-center">
      <div id="sign-in"></div>
    </div>
  `;
  render(template, $.app());
  const signInDiv = document.getElementById("sign-in");
  clerk.mountSignIn(signInDiv);
}

function renderLoading() {
  const template = html`
    <div class="w-full h-full flex items-center justify-center">
      <tree-spinner></tree-spinner>
    </div>
  `;

  render(template, $.app());
}

function renderLobbyPlayerList(players) {
  const lobbyPlayerList = $.lobbyPlayerList();
  if (!lobbyPlayerList) {
    return;
  }

  render(
    html`${players.map(
      (p) =>
        html`<div class="p-2 flex gap-2 items-center">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke-width="1.5"
            stroke="currentColor"
            class="size-4"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M17.982 18.725A7.488 7.488 0 0 0 12 15.75a7.488 7.488 0 0 0-5.982 2.975m11.963 0a9 9 0 1 0-11.963 0m11.963 0A8.966 8.966 0 0 1 12 21a8.966 8.966 0 0 1-5.982-2.275M15 9.75a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z"
            />
          </svg>
          ${p}
        </div>`,
    )}`,
    lobbyPlayerList,
  );
}

async function lobbyEffects(players) {
  //const conn = new WebSocket("");
  //conn.onmessage = function (evt) {
  //  var messages = evt.data.split("\n");
  //  for (var i = 0; i < messages.length; i++) {
  //    players.push(messages[i]);
  //    renderLobbyPlayerList(players);
  //  }
  //};
  //
  //return () => {
  //  conn.close();
  //};
}

async function renderLobby(players) {
  const template = html`
    <div class="w-full h-full flex pt-32 flex-col items-center gap-8 relative ">
      <h1 class="font-bold text-4xl  mx-auto">Lobby</h1>
      <div
        class="flex flex-col divide-neutral-300 divide-y w-48"
        id="lobby-player-list"
      ></div>
      <div class="flex flex-col gap-2 w-48">
        <div
          class="flex items-center justify-between p-2 rounded bg-neutral-100 text-gray-600 cursor-copy"
        >
          <span class="text-sm font-mono" id="connection-string"></span>
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke-width="1.5"
            stroke="currentColor"
            class="size-4"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M15.666 3.888A2.25 2.25 0 0 0 13.5 2.25h-3c-1.03 0-1.9.693-2.166 1.638m7.332 0c.055.194.084.4.084.612v0a.75.75 0 0 1-.75.75H9a.75.75 0 0 1-.75-.75v0c0-.212.03-.418.084-.612m7.332 0c.646.049 1.288.11 1.927.184 1.1.128 1.907 1.077 1.907 2.185V19.5a2.25 2.25 0 0 1-2.25 2.25H6.75A2.25 2.25 0 0 1 4.5 19.5V6.257c0-1.108.806-2.057 1.907-2.185a48.208 48.208 0 0 1 1.927-.184"
            />
          </svg>
        </div>
        <button
          class="p-2 hover:bg-neutral-800 rounded mx-auto bg-black text-neutral-50 w-full"
        >
          Start
        </button>
      </div>
    </div>
  `;

  render(template, $.app());

  renderLobbyPlayerList(players);

  try {
    const res = await fetch("/api/lobby/new", {
      headers: {
        Authorization: `Bearer ${await $.clerk.session.getToken()}`,
      },
    });
    const lobby = await res.json();
    console.log(lobby.ConnString);

    if ($.connectionString) {
      $.connectionString.textContent = lobby.ConnString;
    }
  } catch (e) {
    console.error("new lobby error", e);
  }
}

async function onNewGameButton(_) {
  await renderApp(APP_STATE_LOBBY);
}

function renderMain() {
  const template = html`
    <div class="w-full h-full flex pt-32 flex-col items-center gap-8 relative">
      <h1 class="font-bold text-4xl  mx-auto">Two towers</h1>
      <div class="w-48 flex flex-col gap-2">
        <button
            class="p-2 hover:bg-neutral-800 rounded mx-auto bg-black text-neutral-50 w-full"
            id="new-game-button"
        >
          New game
        </button>
        <div class="w-full">
                <input type="text" class="w-full p-2 bg-neutral-100 rounded-t focus:outline outline-black"></input>
          <button
            class="p-2 hover:bg-neutral-800 rounded-b mx-auto bg-black text-neutral-50 w-full"
          >
            Join game
          </button>
        </div>
      </div>
      <div class="absolute top-0 right-0 p-2">
        <div id="user-button"></div>
      </div>
    </div>
  `;

  render(template, $.app());

  const userBtn = $.userBtn();
  if (!userBtn) {
    console.error("user button not found");
  }
  $.clerk.mountUserButton(userBtn);

  const newGameBtn = $.newGameBtn();
  if (!newGameBtn) {
    console.error("new game button not found");
  }

  newGameBtn.addEventListener("mousedown", onNewGameButton);
}

async function renderApp(state) {
  if ($.cancelEffects) {
    $.cancelEffects();
  }

  switch (state) {
    case APP_STATE_LOADING:
      renderLoading();
      break;
    case APP_STATE_LOGIN:
      const clerk = $.clerk;
      if (!clerk) {
        renderApp(APP_STATE_ERR);
        return;
      }

      if (clerk.user) {
        renderApp(APP_STATE_MAIN);
        return;
      }

      renderLogin(clerk);
      break;
    case APP_STATE_MAIN:
      renderMain();
      break;
    case APP_STATE_LOBBY:
      const players = [];
      renderLobby(players);
      $.cancelEffects = await lobbyEffects(players);
      break;
    case APP_STATE_ERR:
      console.log("rendering error");
      break;
  }
}

window.onload = async function () {
  await renderApp(APP_STATE_LOADING);
  $.clerk = new Clerk(clerkPubKey);
  await $.clerk.load({});

  await renderApp(APP_STATE_LOGIN);
};
