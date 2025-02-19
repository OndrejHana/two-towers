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
  clerk: null,
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

function renderMain() {
  const template = html`
    <div class="w-full h-full flex pt-32 flex-col items-center gap-8 relative">
      <h1 class="font-bold text-4xl  mx-auto">Two towers</h1>
      <div class="w-48 flex flex-col gap-2">
        <button
          class="p-2 hover:bg-neutral-800 rounded mx-auto bg-black text-neutral-50 w-full"
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

  const userBtn = document.getElementById("user-button");
  $.clerk.mountUserButton(userBtn);
}

function renderApp(state) {
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
    case APP_STATE_ERR:
      console.log("rendering error");
      break;
  }
}

async function main() {
  renderApp(APP_STATE_LOADING);
  $.clerk = new Clerk(clerkPubKey);
  await $.clerk.load({});

  renderApp(APP_STATE_LOGIN);
}

main();
