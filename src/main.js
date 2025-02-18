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

function clear() {}

function login(clerk) {
  const template = html`
    <div class="w-full h-full flex items-center justify-center">
      <div id="sign-in"></div>
    </div>
  `;
  render(template, $.app());
  const signInDiv = document.getElementById("sign-in");
  clerk.mountSignIn(signInDiv);
}

function loading() {
  const template = html`
    <div class="w-full h-full flex items-center justify-center">
      <tree-spinner></tree-spinner>
    </div>
  `;

  render(template, $.app());
}

function renderApp(state) {
  clear();

  switch (state) {
    case APP_STATE_LOADING:
      loading();
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

      login(clerk);
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
